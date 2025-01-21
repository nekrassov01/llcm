package main

import (
	_ "embed"
	"fmt"
)

var (
	//go:embed completions/llcm.bash
	bashScript string

	//go:embed completions/llcm.zsh
	zshScript string

	//go:embed completions/llcm.ps1
	pwshScript string
)

type completionShell int

const (
	bash completionShell = iota
	zsh
	pwsh
)

func (t completionShell) String() string {
	switch t {
	case bash:
		return "bash"
	case zsh:
		return "zsh"
	case pwsh:
		return "pwsh"
	default:
		return ""
	}
}

func parseShell(s string) (completionShell, error) {
	switch s {
	case bash.String():
		return bash, nil
	case zsh.String():
		return zsh, nil
	case pwsh.String():
		return pwsh, nil
	default:
		return 0, fmt.Errorf("unsupported shell: %q", s)
	}
}
