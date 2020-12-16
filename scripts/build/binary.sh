#!/usr/bin/env bash
#
# Build a static binary for the host OS/ARCH
#

# windows deps for build
go get github.com/inconshreveable/mousetrap
go get github.com/konsorten/go-windows-terminal-sequences

source ./scripts/build/.variables

echo "Building statically linked certonid"
export CGO_ENABLED=0
# cli
gox -os="darwin windows" -arch="amd64" -output="build/certonid.{{.OS}}.{{.Arch}}" -ldflags "${LDFLAGS}" -verbose ./cli
gox -os="linux" -arch="amd64 arm" -output="build/certonid.{{.OS}}.{{.Arch}}" -ldflags "${LDFLAGS}" -verbose ./cli
# serverless
gox -os="linux" -arch="amd64" -output="build/serverless.{{.OS}}.{{.Arch}}" -ldflags "${LDFLAGS}" -verbose ./serverless
