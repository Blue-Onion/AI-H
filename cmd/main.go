package main

import (
	"context"
	"github.com/Blue-Onion/AI-h/internal/cli"
	"log"
	"os"
	"os/signal"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	if err := cli.Run(ctx, os.Args[1:], os.Stdout, os.Stderr); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
