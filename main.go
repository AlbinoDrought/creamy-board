package main

import (
	"context"
	"os"

	"go.albinodrought.com/creamy-board/internal/cmd"
)

func main() {
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		os.Exit(1)
	}
}
