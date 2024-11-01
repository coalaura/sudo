package main

import (
	"fmt"
	"os"
	"os/exec"
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

	app := shell()
	arg := fmt.Sprintf("/K cd /d %s", cwd)

	if wt, err := exec.LookPath("wt"); err == nil {
		arg = app + " " + arg
		app = wt
	}

	if args[0] == "su" {
		// Don't need to do anything here if we are already elevated
		if isAdmin {
			return
		}
	} else {
		arg += " && " + strings.Join(args, " ") + " && pause && exit"
	}

	if isAdmin {
		app = args[0]
		arg = strings.Join(args[1:], " ")

		regular(app, arg, cwd)
	} else {
		elevate(app, arg, cwd)
	}
}
