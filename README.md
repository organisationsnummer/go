# organisationsnummer [![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/organisationsnummer/go/go.yml?branch=main)](https://github.com/organisationsnummer/go/actions) [![Go Reference](https://pkg.go.dev/badge/github.com/organisationsnummer/go.svg)](https://pkg.go.dev/github.com/organisationsnummer/go) [![Go Report Card](https://goreportcard.com/badge/github.com/organisationsnummer/go)](https://goreportcard.com/report/github.com/organisationsnummer/go)

Validate Swedish organization numbers. Follows version 1.1 of the [specification](https://github.com/organisationsnummer/meta#package-specification-v11).

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

## In memoriam

Fredrik "Frozzare" Forsmo (1991-2026) was the initiator, co-founder and a core contributor of the organisationsnummer project. This library carries his work. He is missed.

## License

MIT
