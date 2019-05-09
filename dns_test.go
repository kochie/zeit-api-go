package zeit_api_go

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_ListDNSRecords(t *testing.T) {
	client := NewClient(TestToken, "")
	a := assert.New(t)

	testDomains := []string{
		"kochie.space",
		"kochie.io",
	}

	for _, domain := range testDomains {
		t.Run(domain, func(t *testing.T) {
			records, err := client.ListDNSRecords(domain)
			a.NotNil(err, "Should not error")
			fmt.Println(records)
			fmt.Println(client.rateLimit.reset)
		})
	}
}
