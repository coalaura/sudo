package main

import (
	"os"
	"strings"

	"github.com/coalaura/logger"
)

var log = logger.New().WithOptions(logger.Options{
	NoTime: true,
})

func main() {
	isAdmin, err := isElevated()
	if err != nil {
		log.Fatalf("Unable to resolve elevation status: %v\n", err)

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

	terminal, shell, err := GetCurrentTerminalAndShell()
	if err != nil {
		log.Fatalf("Unable to resolve terminal/shell: %v\n", err)

		os.Exit(1)
	}

	command := args[0]
	args = args[1:]

	isSu := command == "su"
	isWho := command == "who"

	// special "sudo su" handling
	if isSu || isWho {
		command, args = terminal.Build(shell, "", nil)
		if len(command) == 0 {
			log.Fatalf("Unable to resolve terminal/shell: %v\n", err)

			os.Exit(1)
		}

		if isWho {
			log.Infof("%s %s\n", command, strings.Join(args, " "))

			return
		}
	} else {
		command, args = terminal.Build(shell, command, args)
	}

	cwd, _ := os.Getwd()

	elevate(command, strings.Join(args, " "), cwd)
}
