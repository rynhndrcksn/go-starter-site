package main

import (
	"errors"
	"maps"
	"testing"
	"time"

	"github.com/rynhndrcksn/go-starter-site/internal/assert"
	"github.com/rynhndrcksn/go-starter-site/internal/testdata"
	"github.com/rynhndrcksn/go-starter-site/ui"
)

func TestHashAssetPath(t *testing.T) {
	// Temporarily replace the global ui.Files with the test files until after the test is done.
	originalFiles := ui.Files
	ui.Files = testdata.TestFiles
	defer func() {
		ui.Files = originalFiles
	}()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		{
			name:    "empty path",
			input:   "",
			want:    "",
			wantErr: errHashAssetPathIsEmpty,
		},
		{
			name:    "whitespace path",
			input:   "   ",
			want:    "",
			wantErr: errHashAssetPathIsEmpty,
		},
		{
			name:    "valid file path",
			input:   "/assets/test.txt",
			want:    "/assets/test.txt?v=adfa3ed371a7115e1e4a54503b5fdbe9",
			wantErr: nil,
		},
		{
			name:    "non-existent file",
			input:   "/assets/nonexistent.txt",
			want:    "/assets/nonexistent.txt",
			wantErr: errHashAssetCantReadFile,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := hashAssetPath(tt.input)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got: %v, want: %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
	}

	path := "/assets/test.txt"
	hash1, _ := hashAssetPath(path)
	hash2, _ := hashAssetPath(path)

	t.Run("consistent hashing", func(t *testing.T) {
		if hash1 != hash2 {
			t.Errorf("got: %v, want: %v", hash1, hash2)
		}
	})
}

func BenchmarkHashAssetPath(b *testing.B) {
	// Temporarily replace the global ui.Files with the test files until after the test is done.
	originalFiles := ui.Files
	ui.Files = testdata.TestFiles
	defer func() {
		ui.Files = originalFiles
	}()
	for i := 0; i < b.N; i++ {
		_, _ = hashAssetPath("/assets/test.txt")
	}
}

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

func BenchmarkHumanDate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		humanDate(time.Date(2024, 01, 1, 12, 0, 0, 0, time.UTC))
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
			wantErr: errPropsKeyValueCountIsZero,
		},
		{
			name:    "Mismatched key value pairs",
			input:   []any{"key1", "value1", "key2"},
			wantMap: nil,
			wantErr: errPropsKeyValueCountMismatch,
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

func BenchmarkProps(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = props([]any{"key1", "value1", "key2", 5, "key3", true})
	}
}
