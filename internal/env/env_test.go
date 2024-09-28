package env

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/rynhndrcksn/go-starter-site/internal/assert"
)

func TestGetBoolOrDefault(t *testing.T) {
	tests := []struct {
		name      string
		valExists bool
		want      bool
	}{
		{
			name:      "gets the correct value from the environment",
			valExists: true,
			want:      true,
		},
		{
			name:      "gets the default value if environment variable isn't set",
			valExists: false,
			want:      false,
		},
		{
			name:      "gets the default value if environment variable isn't valid",
			valExists: true,
			want:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.valExists {
				if tt.want == true {
					err := os.Setenv("val", strconv.FormatBool(tt.want))
					if err != nil {
						t.Fatal(err.Error())
					}
				} else {
					err := os.Setenv("val", "not a bool")
					if err != nil {
						t.Fatal(err.Error())
					}
				}
			}

			got := GetBoolOrDefault("val", tt.want)
			assert.Equal(t, got, tt.want)

			if tt.valExists {
				err := os.Unsetenv("val")
				if err != nil {
					t.Fatal(err.Error())
				}
			}
		})
	}
}

func TestGetDurationOrDefault(t *testing.T) {
	tests := []struct {
		name      string
		valExists bool
		want      time.Duration
	}{
		{
			name:      "gets the correct value from the environment",
			valExists: true,
			want:      time.Second,
		},
		{
			name:      "gets the default value if environment variable isn't set",
			valExists: false,
			want:      time.Millisecond,
		},
		{
			name:      "gets the default value if environment variable isn't valid",
			valExists: true,
			want:      time.Millisecond,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.valExists {
				if tt.want == time.Second {
					err := os.Setenv("val", "1s")
					if err != nil {
						t.Fatal(err.Error())
					}
				} else {
					err := os.Setenv("val", "not a time.Duration")
					if err != nil {
						t.Fatal(err.Error())
					}
				}
			}

			got := GetDurationOrDefault("val", tt.want)
			assert.Equal(t, got, tt.want)

			if tt.valExists {
				err := os.Unsetenv("val")
				if err != nil {
					t.Fatal(err.Error())
				}
			}
		})
	}
}

func TestGetIntOrDefault(t *testing.T) {
	tests := []struct {
		name      string
		valExists bool
		want      int
	}{
		{
			name:      "gets the correct value from the environment",
			valExists: true,
			want:      10,
		},
		{
			name:      "gets the default value if environment variable isn't set",
			valExists: false,
			want:      1,
		},
		{
			name:      "gets the default value if environment variable isn't valid",
			valExists: true,
			want:      1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.valExists {
				if tt.want == 10 {
					err := os.Setenv("val", strconv.FormatInt(int64(tt.want), 10))
					if err != nil {
						t.Fatal(err.Error())
					}
				} else {
					err := os.Setenv("val", "not an int")
					if err != nil {
						t.Fatal(err.Error())
					}
				}
			}

			got := GetIntOrDefault("val", tt.want)
			assert.Equal(t, got, tt.want)

			if tt.valExists {
				err := os.Unsetenv("val")
				if err != nil {
					t.Fatal(err.Error())
				}
			}
		})
	}
}

func TestGetStringOrDefault(t *testing.T) {
	tests := []struct {
		name      string
		valExists bool
		want      string
	}{
		{
			name:      "gets the correct value from the environment",
			valExists: true,
			want:      "a string",
		},
		{
			name:      "gets the default value if environment variable isn't set",
			valExists: false,
			want:      "another string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.valExists {
				err := os.Setenv("val", tt.want)
				if err != nil {
					t.Fatal(err.Error())
				}

			}

			got := GetStringOrDefault("val", tt.want)
			assert.Equal(t, got, tt.want)

			if tt.valExists {
				err := os.Unsetenv("val")
				if err != nil {
					t.Fatal(err.Error())
				}
			}
		})
	}
}
