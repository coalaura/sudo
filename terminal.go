package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/shirou/gopsutil/process"
)

type Terminal struct {
	Path string
}

var KnownEnvironments = map[string]string{
	"WT_SESSION":       "wt.exe",
	"WEZTERM_PANE":     "wezterm-gui.exe",
	"ALACRITTY_LOG":    "alacritty.exe",
	"TABBY_PROFILE":    "Tabby.exe",
	"TERMINUS_PLUGINS": "Terminus.exe",
	"ConEmuPID":        "ConEmu64.exe",
	"CMDER_ROOT":       "Cmder.exe",
}

func (t *Terminal) Build(sh *Shell, cmd string, args []string) (string, []string) {
	command := sh.Build(cmd, args)

	if t.Path == "" {
		return command[0], command[1:]
	}

	exe := strings.ToLower(filepath.Base(t.Path))

	switch exe {
	case "wt.exe", "windowsterminal.exe":
		return "wt.exe", append([]string{"new-tab"}, command...)
	case "rio.exe":
		return t.Path, append([]string{"-e"}, command...)
	case "tabby.exe":
		return t.Path, append([]string{"run"}, command...)
	case "alacritty.exe":
		return t.Path, append([]string{"-e"}, command...)
	case "wezterm-gui.exe", "wezterm.exe":
		return t.Path, append([]string{"start", "--"}, command...)
	case "conemu64.exe", "conemu.exe":
		return t.Path, append([]string{"-run"}, command...)
	}

	return command[0], command[1:]
}

func GetCurrentTerminalAndShell() (*Terminal, *Shell, error) {
	// resolve shell first
	pid := int32(os.Getppid())

	shellProc, err := process.NewProcess(pid)
	if err != nil {
		return nil, nil, err
	}

	path, err := shellProc.Exe()
	if err != nil {
		return nil, nil, err
	}

	shell := Shell{path}

	// find terminal from env
	var terminal Terminal

	for env, exe := range KnownEnvironments {
		if os.Getenv(env) == "" {
			continue
		}

		path, err := exec.LookPath(exe)
		if err != nil && exe == "ConEmu64.exe" {
			path, _ = exec.LookPath("ConEmu.exe")
		}

		terminal.Path = path

		break
	}

	// fallback to parent
	if terminal.Path == "" {
		terminalProc, err := shellProc.Parent()
		if err != nil {
			return &terminal, &shell, nil
		}

		path, err = terminalProc.Exe()
		if err != nil {
			return nil, nil, err
		}

		terminal.Path = path
	}

	return &terminal, &shell, nil
}
