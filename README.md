# ZEIT API
[![Build Status](https://travis-ci.com/kochie/zeit-api-go.svg?token=DyduaqJxsshHLt3JzTx3&branch=master)](https://travis-ci.com/kochie/zeit-api-go)
[![GoDoc](https://godoc.org/github.com/kochie/zeit-api-go?status.svg)](https://godoc.org/github.com/kochie/zeit-api-go)
[![Coverage Status](https://coveralls.io/repos/github/kochie/zeit-api-go/badge.svg?branch=master)](https://coveralls.io/github/kochie/zeit-api-go?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/kochie/zeit-api-go)](https://goreportcard.com/report/github.com/kochie/zeit-api-go)
```bash
go get github.com/kochie/zeit-api-go
```

## Example
```go
package main

import (
	"fmt"
	"github.com/kochie/zeit-api-go"
)

func main(){
	token := "secret token"
	
	zeitClient := zeit.NewClient(token)
	zeitClient.Team("team name") // Team name can be optionally set
	
	domains, err := zeitClient.GetAllDomains()
	if err != nil {
		fmt.Println(err.Error())
	}
	
	for _, domain := range domains {
		fmt.Println(domain)
	}
}
```

## Testing
Each method should have sufficient test coverage and an integration test. To facillitate the development of tests there is a mocking interface set up which can be used. Mocks can be created using the `go generate` command or by using the `mockgen` command.
```bash
# Example mocking of HttpClient
mockgen -destination=mocks/mock_http_client.go -package=mocks github.com/kochie/zeit-api-go HttpClient
```

Integration tests will not run unless the `integration` flag is set.

```bash
go test -tags integration
```

More information about setting up a development environment can be found in the [contribution guide](./CONTRIBUTING.md).

## APIs
As listed in the [API documentation](https://zeit.co/docs/api)

Currently supported endpoints are.
- [x] Domains
- [x] DNS
- [ ] OAuth2
- [ ] Authentication
- [ ] Deployments
- [ ] Logs
- [ ] Certificates
- [ ] Aliases
- [ ] Secrets
- [ ] Teams
- [ ] Projects
