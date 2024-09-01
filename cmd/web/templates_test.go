package main

import (
	"maps"
	"testing"
	"time"

	"github.com/rynhndrcksn/go-starter-site/internal/assert"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2024, 01, 1, 12, 0, 0, 0, time.UTC),
			want: "01 Jan 2024 at 12:00",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "PST",
			tm:   time.Date(2024, 01, 1, 12, 0, 0, 0, time.FixedZone("UTC-8", 8*60*60)),
			want: "01 Jan 2024 at 04:00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)
			assert.Equal(t, hd, tt.want)
		})
	}
}

func TestProps(t *testing.T) {
	validMap := make(map[string]any, 2)
	validMap["key1"] = "value1"
	validMap["key2"] = 5
	validMap["key3"] = true
	tests := []struct {
		name    string
		input   []any
		wantMap map[string]any
		wantErr error
	}{
		{
			name:    "Passed in 0 items",
			input:   []any{},
			wantMap: nil,
			wantErr: propsPairsCountIsZeroError,
		},
		{
			name:    "Mismatched key value pairs",
			input:   []any{"key1", "value1", "key2"},
			wantMap: nil,
			wantErr: propsKeyValueCountMismatchedError,
		},
		{
			name:    "Valid amount of key value pairs",
			input:   []any{"key1", "value1", "key2", 5, "key3", true},
			wantMap: validMap,
			wantErr: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotMap, gotErr := props(test.input...)
			assert.Equal(t, gotErr, test.wantErr)
			assert.Equal(t, maps.Equal(gotMap, test.wantMap), true)
		})
	}
}
