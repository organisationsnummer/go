# organisationsnummer [![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/organisationsnummer/go/go.yml?branch=master)](https://github.com/organisationsnummer/go/actions) [![GoDoc](https://godoc.org/github.com/organisationsnummer/go?status.svg)](https://godoc.org/github.com/organisationsnummer/go) [![Go Report Card](https://goreportcard.com/badge/github.com/organisationsnummer/go)](https://goreportcard.com/report/github.com/organisationsnummer/go)

Validate Swedish organization numbers.

## Installation

```
go get -u github.com/organisationsnummer/go
```

## Example

```go
package main

import (
	organisationsnummer "github.com/organisationsnummer/go"
)

func main() {
	organisationsnummer.Valid("202100-5489")
	//=> true
}
```

## License

MIT
