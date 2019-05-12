//+build integration

package zeit

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestClient_Integration_ListDNSRecords(t *testing.T) {
	a := assert.New(t)
	client := NewClient(os.Getenv("NOW_TEST_TOKEN"))
	domains := []string{
		"test.com",
	}
	for _, domain := range domains {
		t.Run(domain, func(t *testing.T) {
			records, err := client.ListDNSRecords(domain)
			a.Nil(err, "Error should be nil")
			a.NotNil(records, "records should be defined")

		})
	}
}

func TestClient_Integration_CreateDNSRecord(t *testing.T) {
	a := assert.New(t)
	client := NewClient(os.Getenv("NOW_TEST_TOKEN"))

	domains := []string{
		"test.com",
	}

	records := []struct{
		name string
		recordType string
		value string
	}{
		{"@", RecordTypeA, "1.1.1.1"},
		{"@", RecordTypeAAAA, "1::1"},
		{"@", RecordTypeALIAS, "google.com"},
		{"@", RecordTypeCAA, "0 issue \"letsencrypt.org\""},
		{"test", RecordTypeCNAME, "google.com"},
		{"@", RecordTypeMX, "aspmx.l.google.com 1"},
		{"_sip._tcp.example.com.", RecordTypeSRV, "10 20 5000 sip-server.example.com."},
		//{"@"}
	}

	for _, domain := range domains {
		for _, record := range records {
			t.Run(fmt.Sprintf("%s_%s", domain, record), func(t *testing.T) {
				uuid, err := client.CreateDNSRecord(domain, record.name, record.recordType, record.value)
				a.Nil(err, "Error should be nil")
				a.NotNil(uuid, "records should be defined")
			})
		}
	}
}

func TestClient_Integration_RemoveDNSRecord(t *testing.T) {

}