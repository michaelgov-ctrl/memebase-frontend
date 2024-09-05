package main

import (
	"net/http"
	"testing"

	"github.com/michaelgov-ctrl/memebase-front/internal/assert"
)

func TestMemeView(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/meme/view/66c78b9d663ce3e2a6bf35c4",
			wantCode: http.StatusOK,
			wantBody: "hilarious meme",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/meme/view/42078b9d663ce3e2a6bf3569",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Invalid ID",
			urlPath:  "/meme/view/error",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/meme/view/",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)
			t.Logf("code: %v; body: %v", code, body)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}
