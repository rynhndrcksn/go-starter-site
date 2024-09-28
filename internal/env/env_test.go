package env

import (
	"os"
	"strconv"
	"testing"

	"github.com/rynhndrcksn/go-starter-site/internal/assert"
)

func TestGetBoolOrDefault(t *testing.T) {
	tests := []struct {
		name      string
		valExists bool
		want      bool
		defaultV  bool
	}{
		{
			name:      "gets the bool value from the environment",
			valExists: true,
			want:      true,
			defaultV:  true,
		},
		{
			name:      "gets the default value if environment variable isn't set",
			valExists: false,
			want:      false,
			defaultV:  false,
		},
		{
			name:      "gets the default value if it isn't a bool",
			valExists: true,
			want:      false,
			defaultV:  false,
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

			got := GetBoolOrDefault("val", tt.defaultV)
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
