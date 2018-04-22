# Changelog

All notable changes to this project will be documented in this file.

## [1.0.0] - 2018-04-22

### Added

- Implement [go-chi](https://github.com/go-chi/chi) as a REST api router.
- Implement [cobra](https://github.com/spf13/cobra) for cli.
- Implement a config handler using [yaml](https://github.com/go-yaml/yaml) encoder and decoder.
- Use [Zap Logger](https://github.com/uber-go/zap) for logging REST api request. And use [lumberjack](https://github.com/natefinch/lumberjack) for log rotation.
- Use golang [dep](https://github.com/golang/dep) for dependencies management.
- Create Makefile for building server binaries.