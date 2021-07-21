#!/bin/sh
set -ex
VERSION=`git describe --abbrev=0 --tags`
dist() {
	export GOOS=${1%/*} GOARCH=${1#*/}
	go build
	tar zcvf .github/otpauth-$VERSION-$GOOS-$GOARCH.tgz LICENSE README.md images/*.png otpauth*
	go clean
}
# see `go tool dist list` for possible target values
dist "linux/amd64"
dist "darwin/amd64"
dist "darwin/arm64"
dist "windows/amd64"
