# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Planned for 0.1.0
- SMTP server based on the standard Go package net/textproto
- REST API for reading mails sent, based on the package github.com/go-chi
- Configuration with defaults settings from environment variables
- Travis CI build configuration
- Makefile
### Planned for 1.0.0
- SMTP server handles correctly edges cases like unexpected connection loss
- OpenAPI and gRPC API description
- Configuration with defaults settings from TOML envfile, environment variables
  or command line flags using package github.com/spf13/viper
- Logging with github.com/goph/logur and github.com/sirupsen/logrus
- Go tests for every file
- Dockerfile and docker-compose.yaml files
- Go releases with https://github.com/goreleaser/goreleaser
- Complete documentation including README.md (https://www.writethedocs.org/)
- Code of conduct (https://www.contributor-covenant.org/)
### Planned for 1.1.0
- Content of mail decoded for example with packages mime or net/mail
- SMTP extensions like authentication, tls
- Addresses decoded for example with package net/mail
- Healthcheck with github.com/InVisionApp/go-health
- Gracefull restarts with github.com/cloudflare/tableflip
- Live reloading of configuration
- Extend API with search service
