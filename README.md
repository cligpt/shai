# shai

[![Build Status](https://github.com/cligpt/shai/workflows/ci/badge.svg?branch=main&event=push)](https://github.com/cligpt/shai/actions?query=workflow%3Aci)
[![codecov](https://codecov.io/gh/cligpt/shai/branch/main/graph/badge.svg?token=El8oiyaIsD)](https://codecov.io/gh/cligpt/shai)
[![Go Report Card](https://goreportcard.com/badge/github.com/cligpt/shai)](https://goreportcard.com/report/github.com/cligpt/shai)
[![License](https://img.shields.io/github/license/cligpt/shai.svg)](https://github.com/cligpt/shai/blob/main/LICENSE)
[![Tag](https://img.shields.io/github/tag/cligpt/shai.svg)](https://github.com/cligpt/shai/tags)



## Introduction

*shai* is a shell with AI written in Go.



## Prerequisites

- Go >= 1.22.0



## Build

```bash
version=latest
make build
```



## Usage

```
shell with ai

Usage:
  shai [flags]

Flags:
  -f, --config-file string   config file (default "$HOME/.shai/shai.yml")
  -h, --help                 help for shai
  -l, --log-level string     log level (DEBUG|INFO|WARN|ERROR) (default "WRAN")
  -v, --version              version for shai
```



## Settings

*shai* parameters can be set in the directory [config](https://github.com/cligpt/shai/blob/main/config).

An example of configuration in [config.yml](https://github.com/cligpt/shai/blob/main/config/config.yml):

```yaml
apiVersion: v1
kind: shai
metadata:
  name: shai
spec:
  drive:
    host: 127.0.0.1
    port: 65050
```



## License

Project License can be found [here](LICENSE).



## Reference

- [bubbletea](https://github.com/charmbracelet/bubbletea)
- [dify-dataset-api](https://docs.dify.ai/guides/knowledge-base/maintain-dataset-via-api)
- [go-basher](https://github.com/progrium/go-basher)
- [ollama-model-api](https://github.com/ollama/ollama/blob/main/docs/api.md)
- [shai](https://github.com/jonboh/shai)
- [survey](https://github.com/AlecAivazis/survey)
- [warp](https://www.warp.dev/)
