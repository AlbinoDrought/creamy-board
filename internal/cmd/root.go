package cmd

import (
	"context"
	"errors"
	"log"
)

type commandHandler func(context.Context, []string) error

var commands map[string]commandHandler = map[string]commandHandler{
	"migrate": migrate,
}

const defaultCommand = "serve"

var ErrCommandNotFound = errors.New("command not found")

func Run(ctx context.Context, args []string) error {
	var command string
	if len(args) > 1 {
		command = args[1]
	} else {
		command = defaultCommand
	}

	handler, ok := commands[command]
	if ok {
		return handler(ctx, args[1:])
	}

	log.Printf("unknown command %v", command)
	return ErrCommandNotFound
}
