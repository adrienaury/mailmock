# Mailmock

[![Go Report Card](https://goreportcard.com/badge/github.com/adrienaury/mailmock)](https://goreportcard.com/report/github.com/adrienaury/mailmock)
[![Github Release Card](https://img.shields.io/github/release/adrienaury/mailmock)](https://github.com/adrienaury/mailmock/releases)
[![codecov](https://codecov.io/gh/adrienaury/mailmock/branch/develop/graph/badge.svg)](https://codecov.io/gh/adrienaury/mailmock)
[![Build Status](https://travis-ci.org/adrienaury/mailmock.svg?branch=develop)](https://travis-ci.org/adrienaury/mailmock)

Mailmock is a lightweight SMTP server designed for testing. It exposes a REST API which will enable your CI/CD to check what transaction were made and what was sent to who.

Mailmock is inspired by Mailhog, but doesn't expose a graphical user interface.

## Features

- SMTP Server implementing RFC5321
- HTTP REST API to list transactions and mails the SMTP server handles

## Installation

### Using Docker

```bash
# this will start mailmock, binding the SMTP server on local port 1025 and the HTTP server to local port 1080
docker run -ti --rm -p 1080:80 -p 1025:25 adrienaury/mailmock
```

### Not using Docker

Download the latest version for your OS from the [release page](https://github.com/adrienaury/mailmock/releases).

Then run it:

```bash
# this will start mailmock with the default configuration,
# binding the SMTP server on local port 25 and the HTTP server to local port 80
mailmock
```

## Configuration

Mailmock can be configured by (in order of precedence) :
- passing flag argument on command line
- setting environment variable
- using a configuration file (JSON, TOML, YAML, HCL, envfile or Java properties formats supported)
- using default provided values

A mix of all of these possibilities can be used.

| Flag argument     | Environment var   | Config file param | Default Value | Description                                                   |
|-------------------|-------------------|-------------------|---------------|---------------------------------------------------------------|
| --logLevel string | MAILMOCK_LOGLEVEL | logLevel          | info          | Set the logger level (trace, debug, info, warn, error)        |
| --httpPort string | MAILMOCK_HTTPPORT | httpPort          | http          | Port number or alias (such as "http") used by the HTTP server |
| --smtpPort string | MAILMOCK_SMTPPORT | smtpPort          | smtp          | Port number or alias (such as "smtp") used by the SMTP server |
| --address string  | MAILMOCK_ADDRESS  | address           |               | IP or hostname                                                |
| --config string   |                   |                   |               | Override default location of configuration file               |

### Configuration file

The configuration file can be placed in different locations :
- /etc/mailmock/
- $HOME/.mailmock/
- ./ (working directory of the mailmock process)
- location given by --config flag if present

It must be named config.ext, possible values for ext : json, toml, yaml, yml, properties, props, prop, hcl, dotenv, env.

#### Examples

- config.yaml
```yaml
logLevel: debug
httpPort: 1234
smtpPort: 4321
address: localhost
```

- config.json
```json
{
    "logLevel": "warn",
    "httpPort": "http",
    "smtpPort": "smtp"
}
```

## Contribute

Contributions to this project are very welcome.

If you want to contribute, please check CONTRIBUTING.md

### Links

- Issue Tracker: github.com/adrienaury/mailmock/issues
- Source Code: github.com/adrienaury/mailmock

## Support

If you are having issues, please let me know.
I'm Adrien and my mail is adrien.aury@gmail.com

## License

### Main license

The project is licensed under the [GNU GENERAL PUBLIC LICENSE v3](https://www.gnu.org/licenses/gpl-3.0.html).

### Exception notices

Some files contains a [GPL linking exception](https://en.wikipedia.org/wiki/GPL_linking_exception) to allow linking modules in any project (everything under pkg folder).

### 3rd party librairies

Library                         | Version             | Licenses                          | Usage               |
--------------------------------|---------------------|-----------------------------------|---------------------|
github.com/logur/adapter-logrus | v0.2.0              | [MIT](NOTICE.md#adapter-logrus)   | Logging             |
github.com/logur/logur          | v0.15.0             | [MIT](NOTICE.md#logur)            | Logging             |
github.com/sirupsen/logrus      | v1.4.2              | [MIT](NOTICE.md#logrus)           | Logging             |
github.com/spf13/viper          | v1.4.0              | [MIT](NOTICE.md#viper)            | Configuration       |
github.com/spf13/pflag          | v1.0.3              | [BSD-3-Clause](NOTICE.md#pflag)   | Configuration       |
github.com/go-chi/chi           | v4.0.2+incompatible | [MIT](NOTICE.md#chi)              | HTTP                |
github.com/go-chi/render        | v1.0.1              | [MIT](NOTICE.md#render)           | HTTP                |
github.com/heptio/workgroup     | v0.8.0-beta.1       | [Apache-2.0](NOTICE.md#workgroup) | Synchronization     |
github.com/stretchr/testify     | v1.3.0              | [MIT](NOTICE.md#testify)          | Testing             |
golang.org/pkg                  | v1.13               | [BSD-3-Clause](NOTICE.md#go)      | Go Standard Library |

Check NOTICE.md for copyright notices.
