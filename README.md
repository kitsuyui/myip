# Usage

```console
$ ./myip
203.0.113.2
```

# Installation

```
$ wget https://github.com/kitsuyui/go-myip/releases/download/0.0.1a/myip_{ your OS } -O myip
$ chmod +x ./myip
```

Choose your OS here: https://github.com/kitsuyui/go-myip/releases

# Build

```console
$ go get github.com/jteeuwen/go-bindata/...
$ go generate
$ go get -d ./...
$ go build
```

## with Docker

```console
$ docker run --rm -v "$(pwd)":/myip -w /myip tcnksm/gox sh -c "./build.sh"
```
