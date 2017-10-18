# multimissile

[![Travis](https://img.shields.io/travis/mercari/multimissile.svg?style=flat-square)](https://travis-ci.org/mercari/multimissile) [![Go Report Card](https://goreportcard.com/badge/github.com/istyle-inc/multimissile)](https://goreportcard.com/report/github.com/istyle-inc/multimissile)

multimissile is [JSON-RPC](http://www.jsonrpc.org/) base API gateway server. It implements [JSON-RPC batch](http://www.jsonrpc.org/specification#batch) endpoints with extended format for HTTP REST requests (see [SPEC](/SPEC.md)). For example, it receives one single JSON-RPC array which defines multiple HTTP requests and converts it into multiple concurrent HTTP requests. If you have multiple backend microservices and need to request them at same time for one transaction, multimissile simplifies it.

# Status

Production ready.

# Requirement

multimissile requires Go1.8 or later.

# Installation

multimissile provides a executable named `msl` to kick server. To install `msl`, use `go get`,

```
$ go get -u github.com/istyle-inc/multimissile/...
```

# Usage

To run `msl`, you must provide configuration path via `-c` option (See [CONFIGURATION.md](/CONFIGURATION.md)) about details and [`config/example.toml`](/config/example.toml) for example usage.

```
$ msl -c config/example.toml
```

Use `-help` to see more options.


# Configuration

See [CONFIGURATION.md](/CONFIGURATION.md) about details.

# Specification

See [SPEC.md](/SPEC.md) about details.

# Committers

 * Tatsuhiko Kubo [@cubicdaiya](https://github.com/cubicdaiya)

# Contribution

Please read the CLA below carefully before submitting your contribution.

https://www.mercari.com/cla/

# License

Copyright 2016 Mercari, Inc.

Licensed under the MIT License.
