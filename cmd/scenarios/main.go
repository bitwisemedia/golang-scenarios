// Command scenarios wires the scenario packages into a runnable binary: a
// small auth-flavoured HTTP server and a spec-file inspector.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bitwisemedia/golang-scenarios/internal/specmerge"
	"github.com/bitwisemedia/golang-scenarios/internal/webauth"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "serve":
		addr := ":8080"
		if len(os.Args) > 2 {
			addr = os.Args[2]
		}
		log.Fatal(serve(addr))
	case "spec":
		if len(os.Args) < 3 {
			usage()
			os.Exit(2)
		}
		f, err := specmerge.Load(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		for _, e := range f.Evals {
			fmt.Printf("%s: %d assertions\n", e.Name, len(e.Assertions))
		}
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s serve [addr] | spec <file>\n", os.Args[0])
}

// serve runs the auth-flavoured demo endpoints.
func serve(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/login", func(w http.ResponseWriter, _ *http.Request) {
		webauth.WriteState(w, "demo-oauth2-state")
		http.Error(w, "redirecting to provider", http.StatusFound)
	})
	mux.HandleFunc("/callback", func(w http.ResponseWriter, _ *http.Request) {
		webauth.WriteSession(w, "demo-session-token")
		p := webauth.Provider{Provider: "demo", Authenticated: true, AutoLogin: true}
		if err := webauth.SetProviderCookie(w, p, 24*time.Hour); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Error(w, "signed in", http.StatusOK)
	})
	mux.HandleFunc("/logout", func(w http.ResponseWriter, _ *http.Request) {
		webauth.ClearSession(w)
		p := webauth.Provider{Provider: "demo"}
		if err := webauth.SetProviderCookie(w, p, 24*time.Hour); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Error(w, "signed out", http.StatusOK)
	})
	server := &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 10 * time.Second}
	return server.ListenAndServe()
}
