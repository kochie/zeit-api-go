//+build integration

package zeit

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

var token = os.Getenv("NOW_TEST_TOKEN")

var domains = []string{
	"kochie.today",
}

var records = []Record{
	{Name: "foo", Type: RecordTypeA, Value: "1.1.1.1"},
	{Name: "foo1", Type: RecordTypeAAAA, Value: "1::1"},
	{Name: "foo2", Type: RecordTypeALIAS, Value: "google.com."},
	{Name: "foo3", Type: RecordTypeCAA, Value: "0 issue \"letsencrypt.org\""},
	{Name: "test", Type: RecordTypeCNAME, Value: "google.com."},
	{Name: "", Type: RecordTypeMX, Value: "aspmx.l.google.com", MxPriority: "10"},
	{Name: "_sip._tcp.example.com", Type: RecordTypeSRV, Value: "20 5000 sip-server.example.com.", SrvPriority: "10"},
	{Name: "", Type: RecordTypeTXT, Value: "Hello"},
}

func init() {
	if token == "" {
		log.Println("The test token is not defined")
	}
}

func TestClient_Integration_ListDNSRecords(t *testing.T) {
	a := assert.New(t)

	client := NewClient(token)

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
	client := NewClient(token)

	for _, domain := range domains {
		for _, record := range records {
			t.Run(fmt.Sprintf("%s_%s_%s_%s", domain, record.Name, record.Type, record.GetValue()), func(t *testing.T) {
				uuid, createError, err := client.CreateDNSRecord(domain, &record)
				a.Nil(err, "Error should be nil")
				a.Nil(createError, "Should not fail to create record")
				a.NotNil(uuid, "records should be defined")

				actualRecords, err := client.ListDNSRecords(domain)
				a.Nil(err, "Getting DNS records should succeed")
				expectedRecord := record
				for _, actualRecord := range actualRecords {
					if actualRecord.Id == uuid {
						//a.Equal(expectedRecord, actualRecord, "records do not match")
						a.Equal(expectedRecord.Name, actualRecord.Name, "Record names should match")
						a.Equal(expectedRecord.Type, actualRecord.Type, "Record type should match")
						a.Equal(expectedRecord.GetValue(), actualRecord.GetValue(), "Record value should match")
					}
				}
			})
		}
	}

}

func TestClient_Integration_RemoveDNSRecord(t *testing.T) {
	a := assert.New(t)
	client := NewClient(token)

	for _, domain := range domains {
		for _, record := range records {
			t.Run(fmt.Sprintf("%s_%s_%s_%s", domain, record.Name, record.Type, record.GetValue()), func(t *testing.T) {
				uuid, createError, err := client.CreateDNSRecord(domain, &record)
				a.Nil(err, "Error should be nil")
				a.Nil(createError, "Should not fail to create record")
				a.NotNil(uuid, "records should be defined")

				err = client.RemoveDNSRecord(domain, uuid)
				a.Nil(err, "Didn't remove record properly")
			})
		}
	}

}
