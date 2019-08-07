# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Planned for 1.0.0
- SMTP server based on the standard Go package net/textproto
- REST API for reading mails sent, based on the package github.com/go-chi
- OpenAPI and gRPC API description
- Configuration with defaults settings, TOML envfile, environment variable
  or command line flags using package github.com/spf13/viper
- Logging with github.com/goph/logur and github.com/sirupsen/logrus
- Go tests
- Dockerfile and docker-compose.yaml files
- Makefile
- Travis CI build configuration
- Complete documentation including README.md (https://www.writethedocs.org/)
- LICENCE file (MIT License)
### Planned for 1.1.0
- Healthcheck with github.com/InVisionApp/go-health
- Gracefull restarts with github.com/cloudflare/tableflip
