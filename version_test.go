package llcm

import (
	"fmt"
	"testing"
)

func TestVersion(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		revision string
		want     string
	}{
		{
			name:     "basic",
			revision: "1234567",
			want:     fmt.Sprintf("%s (revision: 1234567)", version),
		},
		{
			name:     "no revision",
			version:  version,
			revision: "",
			want:     version,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			revision = tt.revision
			if got := Version(); got != tt.want {
				t.Errorf("Version() = %v, want %v", got, tt.want)
			}
		})
	}
}
