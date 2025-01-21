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
			name:    "completion bash",
			args:    []string{name, "completion", bash.String()},
			wantErr: false,
		},
		{
			name:    "completion zsh",
			args:    []string{name, "completion", zsh.String()},
			wantErr: false,
		},
		{
			name:    "completion pwsh",
			args:    []string{name, "completion", pwsh.String()},
			wantErr: false,
		},
		{
			name:    "completion unknown",
			args:    []string{name, "completion", "fish"},
			wantErr: true,
		},
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
			err := newApp(io.Discard, io.Discard).RunContext(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
