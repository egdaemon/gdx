package main

import (
	"context"
	"log"

	"eg/compute/diagx"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LUTC | log.Ltime)
	ctx, done := context.WithTimeout(context.Background(), egenv.TTL())
	defer done()

	c1 := eg.Container("eg.meta.ubuntu.24.04")

	err := eg.Perform(
		ctx,
		eggit.AutoClone,
		eg.Build(c1.BuildFromFile(".eg/Containerfile")),
		eg.Module(ctx, c1, diagx.GenerateSchema),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
