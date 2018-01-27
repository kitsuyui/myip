#!/bin/sh
go get -u github.com/jteeuwen/go-bindata/...
go generate
go get -d ./...
CGO_ENABLE=0 \
gox \
-osarch='darwin/386 darwin/amd64' \
-osarch='linux/386 linux/amd64 linux/arm' \
-osarch='freebsd/386 freebsd/amd64 freebsd/arm' \
-osarch='openbsd/386 openbsd/amd64' \
-osarch='windows/386 windows/amd64' \
-osarch='netbsd/386 netbsd/amd64' \
-ldflags '-w -s' \
-ldflags '-X main.version='"$BUILD_VERSION" \
-output='build/myip_{{.OS}}_{{.Arch}}'
