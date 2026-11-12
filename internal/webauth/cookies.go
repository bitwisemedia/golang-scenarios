// Package webauth handles the browser-facing auth cookies: the session
// cookie, the OAuth CSRF state cookie, and the SPA-visible JSON cookies the
// status page reads to render its sign-in surface.
package webauth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Cookie names. The session and OAuth state cookies carry secrets and are
// written HttpOnly + Secure; the provider cookie is SPA-visible by design and
// must never carry a secret.
const (
	sessionCookie  = "scenarios-auth"
	stateCookie    = "scenarios-oauth2-state"
	providerCookie = "scenarios-auth-provider"
)

const (
	sessionTTL = 8 * time.Hour
	stateTTL   = 5 * time.Minute
)

// WriteSession writes the session cookie.
func WriteSession(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    token,
		Path:     "/",
		MaxAge:   int(sessionTTL.Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

// WriteState writes the OAuth CSRF state cookie ahead of the authorize
// redirect.
func WriteState(w http.ResponseWriter, state string) {
	http.SetCookie(w, &http.Cookie{
		Name:     stateCookie,
		Value:    state,
		Path:     "/",
		MaxAge:   int(stateTTL.Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

// ClearSession expires the session and OAuth state cookies on logout.
func ClearSession(w http.ResponseWriter) {
	clearCookie(w, sessionCookie)
	clearCookie(w, stateCookie)
}

// clearCookie expires one cookie.
func clearCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{Name: name, Value: "", Path: "/", MaxAge: -1, HttpOnly: true, Secure: true})
}

// Provider describes the sign-in surface the SPA renders.
type Provider struct {
	Provider      string `json:"provider"`
	Authenticated bool   `json:"authenticated"`
	AutoLogin     bool   `json:"autoLogin"`
}

// SetProviderCookie writes the SPA-visible base64url JSON provider cookie.
// Not HttpOnly by design — the SPA reads it through document.cookie to render
// the sign-in surface; it must never carry a secret.
func SetProviderCookie(w http.ResponseWriter, p Provider, maxAge time.Duration) error {
	raw, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("cookie %s: %w", providerCookie, err)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     providerCookie,
		Value:    base64.RawURLEncoding.EncodeToString(raw),
		Path:     "/",
		MaxAge:   int(maxAge.Seconds()),
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	return nil
}
