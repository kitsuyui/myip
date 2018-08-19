#!/bin/sh
go get -d ./...
CGO_ENABLE=0 \
gox \
-ldflags '-w -s' \
-ldflags '-X main.version='"$BUILD_VERSION" \
-output='build/myip_{{.OS}}_{{.Arch}}'
