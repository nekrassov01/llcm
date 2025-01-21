package main

import "testing"

func Test_shell_String(t *testing.T) {
	tests := []struct {
		name string
		tr   completionShell
		want string
	}{
		{
			name: "bash",
			tr:   bash,
			want: "bash",
		},
		{
			name: "zsh",
			tr:   zsh,
			want: "zsh",
		},
		{
			name: "pwsh",
			tr:   pwsh,
			want: "pwsh",
		},
		{
			name: "default",
			tr:   3,
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.String(); got != tt.want {
				t.Errorf("completionShell.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseShell(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    completionShell
		wantErr bool
	}{
		{
			name: "bash",
			args: args{
				s: "bash",
			},
			want:    bash,
			wantErr: false,
		},
		{
			name: "zsh",
			args: args{
				s: "zsh",
			},
			want:    zsh,
			wantErr: false,
		},
		{
			name: "pwsh",
			args: args{
				s: "pwsh",
			},
			want:    pwsh,
			wantErr: false,
		},
		{
			name: "unsupported",
			args: args{
				s: "unsupported",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseShell(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseShell() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseShell() = %v, want %v", got, tt.want)
			}
		})
	}
}
