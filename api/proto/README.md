# Protobuf definitions

## Generate stubs

### Go

```bash
# From repo root
mkdir -p api/proto/gen/go
protoc --go_out=api/proto/gen/go --go_opt=paths=source_relative \
  -I api/proto api/proto/telemetry.proto
```

Requires: `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`

### TypeScript (Node)

```bash
npm install -g ts-proto
mkdir -p api/proto/gen/ts
protoc --ts_proto_out=api/proto/gen/ts -I api/proto api/proto/telemetry.proto
```

Consume from `services/persistence` via path or published package.

## Idempotency

`SensorTelemetry.idempotency_key` must be unique per logical event (e.g. UUID + nanosecond timestamp). The persistence layer uses it for UPSERT to avoid duplicates on retries.
