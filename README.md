<p align="center"><a href="#readme"><img src="https://gh.kaos.st/rds-payload-generator.svg"/></a></p>

<p align="center">
  <a href="https://travis-ci.com/essentialkaos/rds-payload-generator"><img src="https://travis-ci.com/essentialkaos/rds-payload-generator.svg"></a>
  <a href="https://goreportcard.com/report/github.com/essentialkaos/rds-payload-generator"><img src="https://goreportcard.com/badge/github.com/essentialkaos/rds-payload-generator"></a>
  <a href="https://codebeat.co/projects/github-com-essentialkaos-rds-payload-generator-master"><img alt="codebeat badge" src="https://codebeat.co/badges/ddab93a0-a00f-4922-8430-09106383ddba" /></a>
  <a href="https://essentialkaos.com/ekol"><img src="https://gh.kaos.st/ekol.svg"></a>
</p>

<br/>

`rds-payload-generator` is simple payload generator for [Redis-Split](https://github.com/essentialkaos/rds).

### Installation

#### From source

Before the initial install allows git to use redirects for [pkg.re](https://github.com/essentialkaos/pkgre) service (reason why you should do this described [here](https://github.com/essentialkaos/pkgre#git-support)):

```
git config --global http.https://pkg.re.followRedirects true
```

To build the RDS Payload Generator from scratch, make sure you have a working Go 1.10+ workspace ([instructions](https://golang.org/doc/install)), then:

```
go get github.com/essentialkaos/rds-payload-generator
```

If you want to update RDS Payload Generator to latest stable release, do:

```
go get -u github.com/essentialkaos/rds-payload-generator
```

#### Prebuilt binaries

You can download prebuilt binaries for Linux and OS X from [EK Apps Repository](https://apps.kaos.st/rds-payload-generator/latest).

### Usage

```
Usage: rds-payload-generator {options}

Options

  --dir, -d dir      Redis-Split main dir
  --keys, -k         Number of keys (10-1000000 default: 5000)
  --ratio, -r        Writes/reads ration (1-100 default: 4)
  --no-color, -nc    Disable colors in output
  --help, -h         Show this help message
  --version, -v      Show version

Examples

  rds-payload-generator -d /srv/redis-split -k 35000 -r 10
  Run tool with custom settings

```

### Build Status

| Branch | Status |
|------------|--------|
| `master` | [![Build Status](https://travis-ci.com/essentialkaos/rds-payload-generator.svg?branch=master)](https://travis-ci.com/essentialkaos/rds-payload-generator) |
| `develop` | [![Build Status](https://travis-ci.com/essentialkaos/rds-payload-generator.svg?branch=develop)](https://travis-ci.com/essentialkaos/rds-payload-generator) |

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/contributing-guidelines#contributing-guidelines).

### License

[EKOL](https://essentialkaos.com/ekol)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
