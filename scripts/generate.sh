#!/usr/bin/env bash

# ensure required commands are installed
command -v sqlc &> /dev/null || { echo "Error: sqlc is required." >&2; exit 1; }
command -v swag &> /dev/null || { echo "Error: swag is required." >&2; exit 1; }
command -v bunx &> /dev/null || { echo "Error: bunx is required." >&2; exit 1; }

# ensure swag is the correct version
if [[ $(swag --version) != *v2* ]]; then
  echo "Error: swag must be version 2"
  exit 1
fi

# ensure running in opendungeon directory
if [[ $PWD != */opendungeon ]]; then
  echo "Error: script must be run in the opendungeon directory" >&2
  exit 1
fi

API_MAIN="cmd/main.go"
SPEC="docs/swagger.yaml"
SCHEMA="web/src/lib/api/schema.d.ts"

# generate database handlers
sqlc generate

# generate openapi spec
swag fmt
swag init --v3.1 --generalInfo $API_MAIN --outputTypes go,yaml

# generate typescript handlers
bunx openapi-typescript $SPEC --output $SCHEMA
