#!/usr/bin/env sh
set -eu

platform="${1:-linux/amd64}"
case "$platform" in
  linux/amd64) tag="transitmanifest-benzhi-amd64:latest" ;;
  linux/arm64) tag="transitmanifest-benzhi-arm64:latest" ;;
  *) echo "unsupported platform: $platform" >&2; exit 2 ;;
esac

docker build \
  --platform "$platform" \
  --build-arg GO_IMAGE=golang:1.23.12 \
  -f benzhi.Dockerfile \
  -t "$tag" \
  .
