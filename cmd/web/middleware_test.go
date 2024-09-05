package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/michaelgov-ctrl/memebase-front/internal/assert"
)

func TestCommonHeaders(t *testing.T) {
	expectedHeaders := []struct {
		header   string
		expected string
	}{
		{
			header:   "Content-Security-Policy",
			expected: "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com",
		},
		{
			header:   "Referrer-Policy",
			expected: "origin-when-cross-origin",
		},
		{
			header:   "X-Content-Type-Options",
			expected: "nosniff",
		},
		{
			header:   "X-Frame-Options",
			expected: "deny",
		},
		{
			header:   "X-XSS-Protection",
			expected: "0",
		},
		{
			header:   "Server",
			expected: "Go",
		},
		{
			header:   "",
			expected: "",
		},
	}

	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	commonHeaders(next).ServeHTTP(rr, r)

	res := rr.Result()

	for _, ht := range expectedHeaders {
		assert.Equal(t, res.Header.Get(ht.header), ht.expected)
	}

	assert.Equal(t, res.StatusCode, http.StatusOK)

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
