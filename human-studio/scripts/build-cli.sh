#!/bin/bash
# Cross-compile the Human CLI binary for all target platforms.
# Run from the human-studio/ directory.
# Outputs binaries to resources/bin/<platform>/

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
OUT_DIR="$(cd "$(dirname "$0")/.." && pwd)/resources/bin"

echo "Building Human CLI binaries from $REPO_ROOT"
echo "Output: $OUT_DIR"
echo ""

mkdir -p "$OUT_DIR"

# Get version from Go module
VERSION=$(cd "$REPO_ROOT" && git describe --tags --always 2>/dev/null || echo "dev")

LDFLAGS="-s -w -X github.com/barun-bash/human/internal/version.Version=$VERSION"

build() {
  local os=$1
  local arch=$2
  local ext=$3
  local outname="human${ext}"
  local outdir="$OUT_DIR/${os}-${arch}"

  echo "  Building ${os}/${arch}..."
  mkdir -p "$outdir"
  GOOS=$os GOARCH=$arch go build -ldflags "$LDFLAGS" -o "$outdir/$outname" "$REPO_ROOT/cmd/human"
  echo "    → $outdir/$outname"
}

# macOS
build darwin amd64 ""
build darwin arm64 ""

# Windows
build windows amd64 ".exe"

# Linux
build linux amd64 ""

echo ""
echo "Done. Binaries in $OUT_DIR"
ls -la "$OUT_DIR"/*/
