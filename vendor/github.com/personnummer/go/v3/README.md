# personnummer [![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/personnummer/go/go.yml?branch=master)](https://github.com/personnummer/go/actions) [![GoDoc](https://godoc.org/github.com/personnummer/go?status.svg)](https://pkg.go.dev/github.com/personnummer/go/v3) [![Go Report Card](https://goreportcard.com/badge/github.com/personnummer/go)](https://goreportcard.com/report/github.com/personnummer/go)

Validate Swedish personal identity numbers.

## Installation

```
go get -u github.com/personnummer/go/v3
```

## Example

```go
package main

import (
	personnummer "github.com/personnummer/go/v3"
)

func main() {
	personnummer.Valid("198507099805")
	//=> true
}
```

## License

MIT
