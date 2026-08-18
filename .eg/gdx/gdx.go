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

// GenerateProto compiles .proto/gdx.capnp into Go types via the capnp
// compiler + capnpc-go plugin. output lands at gdxapi/gdx.capnp.go per
// the schema's $Go.import annotation.
func GenerateProto(ctx context.Context, _ eg.Op) error {
	gruntime := shellruntime()
	return shell.Run(
		ctx,
		gruntime.New("go work sync && go work vendor"),
		// debug: inspect the container environment to see why go list -m
		// can't resolve the capnp module's directory here.
		gruntime.New("go env GOPATH GOMODCACHE GOFLAGS GOWORK"),
		gruntime.New("go list -mod=readonly -m -f '{{.Dir}}' capnproto.org/go/capnp/v3 || true"),
		gruntime.New("find $(go env GOMODCACHE)/capnproto.org -maxdepth 2 2>&1 || true"),
		// vendoring drops non-Go files (std/*.capnp), so resolve the module cache
		// path instead; -mod=readonly (allowed in workspace mode, unlike -mod=mod)
		// skips the vendor dir so .Dir actually resolves.
		gruntime.New("capnp compile -I$(go list -mod=readonly -m -f '{{.Dir}}' capnproto.org/go/capnp/v3)/std --src-prefix=.proto -ogo:. .proto/gdx.capnp"),
	)
}

func Compile() eg.OpFn {
	return eggolang.AutoCompile(
		eggolang.CompileOption.BuildOptions(
			eggolang.Build(
				eggolang.BuildOption.WorkingDirectory(rootdir()),
			),
		),
	)
}

func Test() eg.OpFn {
	return eg.Sequential(
		eggolang.AutoTest(
			eggolang.TestOption.BuildOptions(
				eggolang.Build(
					eggolang.BuildOption.WorkingDirectory(rootdir()),
				),
			),
		),
		eggolang.RecordCoverage,
	)
}
