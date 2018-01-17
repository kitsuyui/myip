# go-myip

Show own Public IP (a.k.a. Global IP; WAN IP; External IP) with reliability by searching multiple way.

[![CircleCI Status](https://circleci.com/gh/kitsuyui/go-myip.svg?style=shield&circle-token=:circle-token)](https://circleci.com/gh/kitsuyui/go-myip)
[![Codecov Status](https://img.shields.io/codecov/c/github/kitsuyui/go-myip.svg)](https://codecov.io/github/kitsuyui/go-myip/)
[![Github All Releases](https://img.shields.io/github/downloads/kitsuyui/go-myip/total.svg)](https://github.com/kitsuyui/go-myip/releases/latest)

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

## LICENSE

### Source

The 3-Clause BSD License. See also LISENCE file.

### statically linked libraries:

- [golang/go](https://github.com/golang/go/) ... [BSD 3-clause "New" or "Revised" License](https://github.com/golang/go/blob/master/LICENSE)
- [miekg/dns](https://github.com/miekg/dns) ... [BSD 3-clause "New" or "Revised" License](https://github.com/miekg/dns/blob/master/LICENSE)
