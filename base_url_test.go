package main

import (
	"net/http/httptest"
	"testing"
)

func TestGetBaseURL(t *testing.T) {
	tests := []struct {
		name      string
		publicURL string
		host      string
		proto     string
		want      string
	}{
		{"forwarded proto", "", "example.com", "https", "https://example.com"},
		{"no proto defaults to https", "", "example.com", "", "https://example.com"},
		{"localhost defaults to http", "", "localhost:8080", "", "http://localhost:8080"},
		{"PUBLIC_URL overrides wrong proxy proto", "https://app.back4app.io/", "app.back4app.io", "http", "https://app.back4app.io"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("PUBLIC_URL", tt.publicURL)
			r := httptest.NewRequest("GET", "/mcp", nil)
			r.Host = tt.host
			if tt.proto != "" {
				r.Header.Set("X-Forwarded-Proto", tt.proto)
			}
			if got := getBaseURL(r); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
