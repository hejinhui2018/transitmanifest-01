#!/usr/bin/env sh
set -eu

root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
output="$root/bin"
mkdir -p "$output"

for arch in amd64 arm64; do
    docker run --rm \
        -v "$root:/src" \
        -w /src \
        -e CGO_ENABLED=0 \
        -e GOOS=linux \
        -e GOARCH="$arch" \
        golang:1.23.12 \
        go build -trimpath -ldflags="-s -w" -o "/src/bin/transitmanifest-linux-$arch" ./cmd/transitmanifest
done

printf '%s\n' "built linux/amd64 and linux/arm64 binaries in $output"

