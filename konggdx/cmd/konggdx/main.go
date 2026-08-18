package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/egdaemon/gdx/konggdx"
	"github.com/egdaemon/gdx/konggdx/userx"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	var cli konggdx.Commands
	kctx := kong.Parse(&cli,
		kong.Bind(ctx),
		kong.Vars{"vars_gdx_socket": userx.RuntimeDirectory("gdx.socket")},
	)
	kctx.FatalIfErrorf(kctx.Run())
}
