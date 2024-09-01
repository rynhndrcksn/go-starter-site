package main

import (
	"net/http"
	"testing"
)

func TestBackground(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		r    *http.Request
		fn   func()
	}{
		{
			name: "Background works normally",
			r:    req,
			fn: func() {
				t.Log("This is a background test")
			},
		},
		{
			name: "Background recovers from panics",
			r:    req,
			fn: func() {
				panic("this is a panic")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := newTestApplication(t)
			app.background(tt.r, tt.fn)
		})
	}
}
