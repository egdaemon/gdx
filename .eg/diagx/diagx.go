package gdx

import (
	"context"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/eggolang"
)

func rootdir() string {
	return egenv.WorkingDirectory()
}

func shellruntime() shell.Command {
	return eggolang.Runtime().Directory(rootdir())
}

// GenerateSchema compiles .proto/diagx.capnp into Go types via the capnp
// compiler + capnpc-go plugin. output lands at diagxapi/diagx.capnp.go per
// the schema's $Go.import annotation.
func GenerateSchema(ctx context.Context, _ eg.Op) error {
	gruntime := shellruntime()
	return shell.Run(
		ctx,
		gruntime.New("capnp compile -I$(go list -m -f '{{.Dir}}' capnproto.org/go/capnp/v3)/std --src-prefix=.proto -ogo:. .proto/diagx.capnp"),
	)
}
