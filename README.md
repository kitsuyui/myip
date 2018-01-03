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

## LICENSE

### Source

The 3-Clause BSD License. See also LISENCE file.

### statically linked libraries:

- [golang/go](https://github.com/golang/go/) ... [BSD 3-clause "New" or "Revised" License](https://github.com/golang/go/blob/master/LICENSE)
- [bitly/go-simplejson](https://github.com/bitly/go-simplejson) ... [MIT License](https://github.com/bitly/go-simplejson/blob/master/LICENSE)
- [miekg/dns](https://github.com/miekg/dns) ... [BSD 3-clause "New" or "Revised" License](https://github.com/miekg/dns/blob/master/LICENSE)
