package cmd

import (
	"context"
	"errors"

	"go.albinodrought.com/creamy-board/internal/log"
)

var ErrCommandNotFound = errors.New("command not found")

func Run(ctx context.Context, args []string) error {
	command := "serve"
	if len(args) > 1 {
		command = args[1]
	}

	if command == "serve" {
		return serve(ctx)
	}
	if command == "migrate" {
		return migrate(ctx)
	}

	log.Printf("unknown command %v", command)
	return ErrCommandNotFound
}
