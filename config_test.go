package llcm

import (
	"bytes"
	"context"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	type args struct {
		ctx     context.Context
		profile string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				profile: "",
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				profile: "invalid-profile",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			_, err := LoadConfig(tt.args.ctx, tt.args.profile)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.want {
				t.Errorf("LoadConfig() = %v, want %v", gotW, tt.want)
			}
		})
	}
}
