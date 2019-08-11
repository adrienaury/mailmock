# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Planned for 0.2.0
- Cleaner logs, color disabled by default but possibility to reactivate
- Adding pagination to the REST API
### Planned for 1.0.0
- SMTP server handles correctly edges cases like unexpected connection loss
- OpenAPI and gRPC API description
- Configuration with defaults settings from TOML envfile, environment variables
  or command line flags using package github.com/spf13/viper
- Logging with github.com/goph/logur and github.com/sirupsen/logrus
- Go tests for every file
- Dockerfile and docker-compose.yaml files
- Complete documentation
- Code of conduct (https://www.contributor-covenant.org/)
- Logo
### Planned for 1.1.0
- Content of mail decoded for example with packages mime or net/mail
- SMTP extensions like authentication, TLS, ...
- Parsing of addresses with package net/mail
- Healthcheck with github.com/InVisionApp/go-health
- Gracefull restarts with github.com/cloudflare/tableflip
- Live reloading of configuration
- Extend API with search service

## [0.1.1] - TODO
### Fixed
- VRFY command takes exactly 1 parameter

## [0.1.0] - 2019-08-10
### Added
- Dummy SMTP server (RFC 5321 compliant)
- REST API to list transactions and mails the SMTP server handled
- Configuration with defaults settings from environment variables
