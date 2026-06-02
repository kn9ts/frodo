package frodo

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestNewSetsDefaultRouterOptions(t *testing.T) {
	r := New()

	if !r.RedirectTrailingSlash {
		t.Fatal("expected RedirectTrailingSlash to be enabled")
	}
	if !r.RedirectFixedPath {
		t.Fatal("expected RedirectFixedPath to be enabled")
	}
	if !r.HandleMethodNotAllowed {
		t.Fatal("expected HandleMethodNotAllowed to be enabled")
	}
}

func TestRouterServeHTTPChainsHandlers(t *testing.T) {
	r := New()
	var calls []string

	r.Get("/chain", func(w http.ResponseWriter, req *Request) {
		calls = append(calls, "first")
		req.Next()
	}, func(w http.ResponseWriter, req *Request) {
		calls = append(calls, "second")
		fmt.Fprint(w, "ok")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://example.com/chain", nil)

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if got := rec.Body.String(); got != "ok" {
		t.Fatalf("expected response body %q, got %q", "ok", got)
	}
	if !reflect.DeepEqual(calls, []string{"first", "second"}) {
		t.Fatalf("unexpected middleware order: %#v", calls)
	}
}

func TestRouterRedirectsTrailingSlash(t *testing.T) {
	r := New()
	r.Get("/chain", func(w http.ResponseWriter, req *Request) {
		fmt.Fprint(w, "ok")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://example.com/chain/", nil)

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusMovedPermanently {
		t.Fatalf("expected status %d, got %d", http.StatusMovedPermanently, rec.Code)
	}
	if got := rec.Header().Get("Location"); got != "/chain" {
		t.Fatalf("expected redirect location %q, got %q", "/chain", got)
	}
}

func TestRouterReturns404ForMissingRoute(t *testing.T) {
	r := New()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://example.com/missing", nil)

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestRequestClientIPPrefersForwardedHeaders(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	req.RemoteAddr = "203.0.113.10:4321"
	req.Header.Set("X-Forwarded-For", "198.51.100.1, 198.51.100.2")
	r := &Request{Request: req}

	if got := r.ClientIP(); got != "198.51.100.1" {
		t.Fatalf("expected forwarded client IP, got %q", got)
	}

	req.Header.Set("X-Real-Ip", "198.51.100.99")
	if got := r.ClientIP(); got != "198.51.100.99" {
		t.Fatalf("expected real client IP, got %q", got)
	}
}

func TestResponseWriterWriteTracksState(t *testing.T) {
	rec := httptest.NewRecorder()
	w := &ResponseWriter{ResponseWriter: rec, timeStart: time.Now()}

	if _, err := w.Write([]byte("hello")); err != nil {
		t.Fatalf("unexpected write error: %v", err)
	}
	if !w.HeaderWritten() {
		t.Fatal("expected header to be marked written")
	}
	if !w.ResponseSent() {
		t.Fatal("expected response to be marked sent")
	}
	if got := rec.Code; got != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, got)
	}
	if got := rec.Body.String(); got != "hello" {
		t.Fatalf("expected body %q, got %q", "hello", got)
	}
	if got := w.Size(); got != 5 {
		t.Fatalf("expected size %d, got %d", 5, got)
	}
}
