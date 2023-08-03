<p align="center"><a href="#readme"><img src="https://gh.kaos.st/rds-payload-generator.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/w/rds-payload-generator/ci"><img src="https://kaos.sh/w/rds-payload-generator/ci.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/r/rds-payload-generator"><img src="https://kaos.sh/r/rds-payload-generator.svg" alt="GoReportCard" /></a>
  <a href="https://kaos.sh/b/rds-payload-generator"><img src="https://kaos.sh/b/ddab93a0-a00f-4922-8430-09106383ddba.svg" alt="Codebeat badge" /></a>
  <a href="#license"><img src="https://gh.kaos.st/apache2.svg"></a>
</p>

<p align="center"><a href="#installation">Installation</a> • <a href="#usage">Usage</a> • <a href="#build-status">Build Status</a> • <a href="#contributing">Contributing</a> • <a href="#license">License</a></p>

<br/>

`rds-payload-generator` is simple payload generator for [RDS](https://kaos.sh/rds).

### Installation

#### From source

To build the RDS Payload Generator from scratch, make sure you have a working Go 1.20+ workspace ([instructions](https://go.dev/doc/install)), then:

```
go install github.com/essentialkaos/rds-payload-generator@latest
```

#### Prebuilt binaries

You can download prebuilt binaries for Linux and macOS from [EK Apps Repository](https://apps.kaos.st/rds-payload-generator/latest):

```bash
bash <(curl -fsSL https://apps.kaos.st/get) rds-payload-generator
```

### Usage

```
Usage: rds-payload-generator {options}

Options

  --dir, -d dir      Path to RDS main dir
  --keys, -k         Number of keys (10-1000000 | default: 5000)
  --ratio, -r        Writes/reads ratio (1-100 | default: 4)
  --pause, -p        Max pause between requests in ms (1-1000 | default: 15)
  --no-color, -nc    Disable colors in output
  --help, -h         Show this help message
  --version, -v      Show version

Examples

  rds-payload-generator -d /srv/rds -k 35000 -r 10
  Run tool with custom settings
```

### Build Status

| Branch | Status |
|--------|--------|
| `master` | [![CI](https://kaos.sh/w/rds-payload-generator/ci.svg?branch=master)](https://kaos.sh/w/rds-payload-generator/ci?query=branch:master) |
| `develop` | [![CI](https://kaos.sh/w/rds-payload-generator/ci.svg?branch=develop)](https://kaos.sh/w/rds-payload-generator/ci?query=branch:develop) |

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/contributing-guidelines#contributing-guidelines).

### License

[Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
