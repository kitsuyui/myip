# myip

Global IP detection tool for who has paranoid skeptical.
Show own Public IP (a.k.a. Global IP; WAN IP; External IP) with reliability by searching multiple way.

[![Github All Releases](https://img.shields.io/github/downloads/kitsuyui/scraper/total.svg)](https://github.com/kitsuyui/scraper/releases/latest)
[![Codecov Status](https://img.shields.io/codecov/c/github/kitsuyui/myip.svg)](https://codecov.io/github/kitsuyui/myip/)
[![Downloads](https://img.shields.io/github/downloads/kitsuyui/myip/total.svg)](https://github.com/kitsuyui/myip/releases/latest)

# Mechanism

myip supports various way to IP resolving.

- HTTP / HTTPS
- DNS
- STUN

c.f. https://github.com/kitsuyui/myip/blob/master/targets/targets.go

And if many of these targets answer same, then it show as true IP.
So this aims to be resistant for man-in-the-middle-attack or temporary service down.

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
v0.3.0
$
```

## Verbose mode

```console
$ myip -v
2018/01/01 12:00:00 IP:203.0.113.2	type:dns	weight:2.0	A,myip.opendns.com.,resolver3.opendns.com:53
2018/01/01 12:00:00 IP:203.0.113.2	type:dns	weight:2.0	A,whoami.akamai.net.,ns1-1.akamaitech.net:53
2018/01/01 12:00:00 IP:203.0.113.2	type:dns	weight:2.0	TXT,o-o.myaddr.l.google.com.,ns1.google.com:53
2018/01/01 12:00:00 IP:203.0.113.2	type:stun	weight:3.0	stun4.l.google.com:19302
2018/01/01 12:00:00 IP:203.0.113.2	type:http	weight:0.5	http://ipcheck.ieserver.net/
.
.
.
2018/01/01 12:00:00 IP:203.0.113.2	type:https	weight:3.0	https://icanhazip.com/
2018/01/01 12:00:00 IP:203.0.113.2	type:https	weight:3.0	https://wtfismyip.com/text
203.0.113.2
```

## Ready for IPv6

When your machine has IPv6 address, `-6` option is available.

```console
$ myip -6
2001:0db8:0000:0042:0000:8a2e:0370:7334
```

If your machine doesn't have IPv6 address, fallbacks to IPv4.

## Help

```
$ myip -h
myip

Usage:
 myip [-v | --verbose] [-4 | -6] [-T=<rate>] [-t=<duration>]
 myip (--help | --version)

Options:
 -h --help               						 Show this screen.
 -V --version            						 Show version.
 -v --verbose            						 Verbose mode.
 -4 --ipv4               						 Prefer IPv4.
 -6 --ipv6               						 Prefer IPv6.
 -n --newline            						 Show IP with newline.
 -N --no-newline         						 Show IP without newline.
 -T=<rate> --threshold=<rate>  			 Threshold that must be exceeded by weighted votes [default: 0.5].
 -t=<duration> --timeout=<duration>  Timeout [default: 3s].
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

1. Choose your OS here: https://github.com/kitsuyui/myip/releases
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
$ go get -d ./...
$ go build
```

## LICENSE

### Source

The 3-Clause BSD License. See also LICENSE file.

### statically linked libraries:

- [golang/go](https://github.com/golang/go/) ... [BSD 3-Clause "New" or "Revised" License](https://github.com/golang/go/blob/master/LICENSE)
- [miekg/dns](https://github.com/miekg/dns) ... [BSD 3-Clause "New" or "Revised" License](https://github.com/miekg/dns/blob/master/LICENSE)
- [gortc/stun](https://github.com/gortc/stun) ... [BSD 3-Clause "New" or "Revised" License](https://github.com/gortc/stun/blob/master/LICENSE)
- [docopt/docopt-go](https://github.com/docopt/docopt.go) ... [MIT License](https://github.com/docopt/docopt.go/blob/master/LICENSE)
