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
	pathToUrls := map[string]string{
		"/foo": "https://example.com",
	}

	fbHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "FALLBACK")
	})

	t.Run("provided url", func(t *testing.T) {
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
	
	t.Run("not provided url", func(t *testing.T) {
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
