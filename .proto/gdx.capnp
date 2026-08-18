@0xdd5f8a3b1c2e9074;

using Go = import "/go.capnp";
$Go.package("gdx");
$Go.import("github.com/egdaemon/gdx");

# ProfileMode enumerates the supported runtime/pprof profile captures.
enum ProfileMode {
  cpu @0;
  heap @1;
  mem @2;
  allocs @3;
  block @4;
}

# ProfileRequest configures a single profile capture. duration is in
# nanoseconds; 0 means "use the server's configured default".
struct ProfileRequest {
  mode @0 :ProfileMode;
  duration @1 :UInt64;
}

# ProfileResponse carries the raw pprof-encoded profile bytes.
struct ProfileResponse {
  data @0 :Data;
}

# TraceRequest configures a runtime/trace capture. duration is in
# nanoseconds; 0 means "use the server's configured default".
struct TraceRequest {
  duration @0 :UInt64;
}

# TraceResponse carries the raw runtime/trace-encoded bytes.
struct TraceResponse {
  data @0 :Data;
}

# GoroutinesResponse carries a goroutine stack dump, in the same text format
# runtime/pprof.Lookup("goroutine").WriteTo(w, 1) produces.
struct GoroutinesResponse {
  data @0 :Data;
}

# ExpvarEntry is a single published expvar.Var. value is that var's own
# String() output, which expvar guarantees is valid JSON.
struct ExpvarEntry {
  key @0 :Text;
  value @1 :Text;
}

# ExpvarResponse carries every published expvar.Var.
struct ExpvarResponse {
  entries @0 :List(ExpvarEntry);
}
