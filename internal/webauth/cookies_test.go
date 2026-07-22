package webauth

import (
	"net/http/httptest"
	"testing"
	"time"
)

func TestWriteSessionHardened(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteSession(rec, "token")

	cookies := rec.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("got %d cookies, want 1", len(cookies))
	}
	c := cookies[0]
	if c.Name != sessionCookie {
		t.Errorf("name = %q, want %q", c.Name, sessionCookie)
	}
	if !c.HttpOnly || !c.Secure {
		t.Errorf("HttpOnly = %v, Secure = %v, want both true", c.HttpOnly, c.Secure)
	}
}

func TestSetProviderCookieSPAVisible(t *testing.T) {
	rec := httptest.NewRecorder()
	p := Provider{Provider: "demo", Authenticated: true}
	if err := SetProviderCookie(rec, p, time.Hour); err != nil {
		t.Fatal(err)
	}

	cookies := rec.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("got %d cookies, want 1", len(cookies))
	}
	c := cookies[0]
	if c.Name != providerCookie {
		t.Errorf("name = %q, want %q", c.Name, providerCookie)
	}
	if c.HttpOnly {
		t.Error("provider cookie is HttpOnly; the SPA must be able to read it")
	}
}
