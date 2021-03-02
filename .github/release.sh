#!/bin/sh
set -ex
dist() {
	export GOOS=${1%/*} GOARCH=${1#*/}
	go build
	tar zcvf .github/otpauth-$GOOS-$GOARCH.tgz LICENSE README.md otpauth*
	go clean
}
# see `go tool dist list` for possible target values
dist "linux/amd64"
dist "darwin/amd64"
dist "darwin/arm64"
dist "windows/amd64"
