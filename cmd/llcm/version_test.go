package main

import (
	"fmt"
	"testing"
)

func Test_version(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		revision string
		want     string
	}{
		{
			name:     "basic",
			revision: "1234567",
			want:     fmt.Sprintf("%s (revision: 1234567)", Version),
		},
		{
			name:     "no revision",
			version:  Version,
			revision: "",
			want:     Version,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Revision = tt.revision
			if got := version(); got != tt.want {
				t.Errorf("version() = %v, want %v", got, tt.want)
			}
		})
	}
}
