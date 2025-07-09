package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Shell struct {
	Path string
}

func (s *Shell) Build(cmd string, args []string) []string {
	if cmd == "" {
		return []string{s.Path}
	}

	var command []string

	exe := strings.ToLower(filepath.Base(s.Path))

	switch exe {
	case "cmd.exe":
		command = []string{s.Path, "/c"}
	case "powershell.exe", "pwsh.exe":
		command = []string{s.Path, "-NoProfile", "-Command"}
	case "bash.exe", "zsh.exe", "fish.exe", "git-bash.exe", "nu.exe":
		command = []string{s.Path, "-c"}
	}

	if len(command) == 0 {
		return append([]string{cmd}, args...)
	}

	if len(args) > 0 {
		cmd = fmt.Sprintf("%s %s", cmd, strings.Join(args, " "))
	}

	return append(command, cmd)
}
