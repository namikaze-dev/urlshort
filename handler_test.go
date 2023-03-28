package urlshort_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gophercises/urlshort"
)

func TestMapHandler(t *testing.T) {
	t.Run("url provided", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/foo", nil)
		respRec := httptest.NewRecorder()

		mapHandler := urlshort.MapHandler(pathToUrls, fbHandler)
		mapHandler.ServeHTTP(respRec, req)
		resp := respRec.Result()
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("unexpected error from io.ReadAll: %v", err)
		}

		got := string(bytes.TrimSpace(body))
		want := "<a href=\"https://example.com\">See Other</a>."
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("url not provided", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/bar", nil)
		respRec := httptest.NewRecorder()

		mapHandler := urlshort.MapHandler(pathToUrls, fbHandler)
		mapHandler.ServeHTTP(respRec, req)
		resp := respRec.Result()
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("unexpected error from io.ReadAll: %v", err)
		}

		got := string(body)
		want := "FALLBACK"
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
}

func TestYAMLHandler(t *testing.T) {
	t.Run("url provided", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/bar", nil)
		respRec := httptest.NewRecorder()

		YAMLHandler, err := urlshort.YAMLHandler([]byte(yaml), fbHandler)
		if err != nil {
			t.Fatalf("unexpected error from urlshort.YAMLHandler: %v", err)
		}

		YAMLHandler.ServeHTTP(respRec, req)
		resp := respRec.Result()
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("unexpected error from io.ReadAll: %v", err)
		}

		got := string(bytes.TrimSpace(body))
		want := "<a href=\"https://google.com\">See Other</a>."
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("url not provided", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/baz", nil)
		respRec := httptest.NewRecorder()

		YAMLHandler, err := urlshort.YAMLHandler([]byte(yaml), fbHandler)
		if err != nil {
			t.Fatalf("unexpected error from urlshort.YAMLHandler: %v", err)
		}

		YAMLHandler.ServeHTTP(respRec, req)
		resp := respRec.Result()
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("unexpected error from io.ReadAll: %v", err)
		}

		got := string(bytes.TrimSpace(body))
		want := "FALLBACK"
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("invalid yaml input", func(t *testing.T) {
		yaml := `{"foo": "bar"}`
		_, err := urlshort.YAMLHandler([]byte(yaml), fbHandler)
		if err == nil {
			t.Fatal("expected error from urlshort.YAMLHandler, got nil")
		}
	})
}

var pathToUrls = map[string]string{
	"/foo": "https://example.com",
}

var yaml = `
- path: /foo
  url: https://github.com
- path: /bar
  url: https://google.com`

var fbHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "FALLBACK")
})
