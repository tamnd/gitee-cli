// Command gitee is a single-binary command line for Gitee (码云).
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/tamnd/any-cli/kit"
	"github.com/tamnd/gitee-cli/cli"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	os.Exit(kit.Run(ctx, cli.NewApp()))
}
