package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/coalaura/logger"
)

var log = logger.New().WithOptions(logger.Options{
	NoTime: true,
})

func main() {
	isAdmin, err := isElevated()
	if err != nil {
		log.FatalF("Unable to resolve elevation status: %v", err)

		os.Exit(1)
	}

	args := os.Args[1:]

	if len(args) == 0 {
		log.Fatal("Usage: sudo <command>")

		os.Exit(1)
	}

	if isAdmin {
		log.Note("Already elevated.")
	}

	cwd, _ := os.Getwd()

	var (
		app  string
		keep bool
	)

	if args[0] == "su" {
		// Don't need to do anything here if we are already elevated
		if isAdmin {
			return
		}

		// --keep or -k to not close current window
		for _, arg := range args[0:] {
			if arg == "-k" || arg == "--keep" {
				keep = true

				break
			}
		}
	} else {
		app = strings.Join(args, " ")
	}

	cmd, arg := build(app, cwd, isAdmin)

	if isAdmin {
		regular(cmd, arg, cwd)
	} else {
		elevate(cmd, arg, cwd)
	}

	if app == "" && !keep {
		terminateParent()
	}
}

func build(app, cwd string, isAdmin bool) (string, string) {
	var cmd string

	if isAdmin {
		cmd = app
	} else {
		if app == "" {
			cmd = fmt.Sprintf("%s /K \"cd /d %s\"", shellWithWT(), cwd)
		} else {
			rgx := regexp.MustCompile(`\s*;\s*`)
			app = rgx.ReplaceAllString(app, " && ")

			var sh string

			if isShell(app) {
				sh = shell()
			} else {
				sh = shellWithWT()
			}

			cmd = fmt.Sprintf("%s /K \"cd /d %s && %s && exit\"", sh, cwd, app)
		}
	}

	parts := strings.Split(cmd, " ")

	return parts[0], strings.Join(parts[1:], " ")
}
