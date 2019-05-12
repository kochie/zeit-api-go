# ZEIT API
[![Build Status](https://travis-ci.com/kochie/zeit-api-go.svg?token=DyduaqJxsshHLt3JzTx3&branch=master)](https://travis-ci.com/kochie/zeit-api-go)

```
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
	
	client := zeit.NewClient(token)
	client.Team("team name")
	
	domains, err :=client.GetAllDomains()
	if err != nil {
		fmt.Println(err.Error())
	}
	
	for _, domain := range domains {
		fmt.Println(domain)
	}
}
```

## Testing
Mocks can be created using the `go generate` command
```
go generate
```

Integration tests will not run unless the `integration` flag is set.

```bash
go test -tags integration
```

## APIs
As listed in the [API documentation](https://zeit.co/docs/api)

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