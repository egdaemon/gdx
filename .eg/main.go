package main

import (
	"context"
	"log"

	"eg/compute/gdx"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LUTC | log.Ltime)
	ctx, done := context.WithTimeout(context.Background(), egenv.TTL())
	defer done()

	c1 := eg.Container("gdx.ubuntu")

	err := eg.Perform(
		ctx,
		eggit.AutoClone,
		eg.Build(c1.BuildFromFile(".eg/Containerfile")),
		eg.Module(
			ctx,
			c1,
			eg.Sequential(
				gdx.Compile(),
				gdx.Test(),
			),
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
