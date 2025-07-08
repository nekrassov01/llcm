package main

import (
	"context"
	"io"
	"testing"
)

func Test_cli(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "unknown profile",
			args:    []string{name, "-p", "unknown"},
			wantErr: true,
		},
		{
			name:    "unknown log level",
			args:    []string{name, "-l", "unknown"},
			wantErr: true,
		},
		{
			name:    "unknown output type",
			args:    []string{name, "list", "-o", "unknown"},
			wantErr: true,
		},
		{
			name:    "unknown region",
			args:    []string{name, "list", "-r", "unknown"},
			wantErr: true,
		},
		{
			name:    "unknown desired state",
			args:    []string{name, "preview", "-d", "unknown"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := newCmd(io.Discard, io.Discard).Run(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
