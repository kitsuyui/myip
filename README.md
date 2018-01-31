# go-myip

Show own Public IP (a.k.a. Global IP; WAN IP; External IP) with reliability by searching multiple way.

[![CircleCI Status](https://circleci.com/gh/kitsuyui/go-myip.svg?style=shield&circle-token=:circle-token)](https://circleci.com/gh/kitsuyui/go-myip)
[![Codecov Status](https://img.shields.io/codecov/c/github/kitsuyui/go-myip.svg)](https://codecov.io/github/kitsuyui/go-myip/)
[![Github All Releases](https://img.shields.io/github/downloads/kitsuyui/go-myip/total.svg)](https://github.com/kitsuyui/go-myip/releases/latest)

# Usage

```console
$ myip
203.0.113.2
```

## Output with newline by (`-n` or `--newline`) option

This option gives output with ending newline character.

```console
$ myip
203.0.113.2$
```

```console
$ myip -n
203.0.113.2
$
```

## Version

```console
$ myip -V
v0.2.2
$
```

# Installation

## go get

If you have Golang environment, install with just doing this;

```console
$ go get github.com/kitsuyui/myip
```

## Install static binary releases

If you don't have Golang environments or prefer single binary environment, you can use statically binary release.
It has no DLL dependency. So you can use it by just downloading.

1. Choose your OS here: https://github.com/kitsuyui/go-myip/releases
2. Download and make it executable

### Example command

```console
$ wget https://github.com/kitsuyui/myip/releases/download/${version}/myip_${your_os} -O myip
$ chmod +x ./myip
```

## Homebrew / macOS

[kitsuyui/homebrew-myip](https://github.com/kitsuyui/homebrew-myip) is available.

```console
$ brew tap kitsuyui/homebrew-myip
$ brew install myip
```

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
