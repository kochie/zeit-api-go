package zeit

import (
	"fmt"
)

const (
	RecordTypeA     = "A"
	RecordTypeAAAA  = "AAAA"
	RecordTypeALIAS = "ALIAS"
	RecordTypeCAA   = "CAA"
	RecordTypeCNAME = "CNAME"
	RecordTypeMX    = "MX"
	RecordTypeSRV   = "SRV"
	RecordTypeTXT   = "TXT"
)

type Record struct {
	Id          string
	Slug        string
	Type        string
	Name        string
	Value       string
	Creator     string
	Created     *Time
	Updated     *Time
	MxPriority  string `json:"mxPriority,omitempty"`
	SrvPriority string `json:"priority,omitempty"`
}

func (r *Record) GetValue() string {
	switch r.Type {
	case RecordTypeSRV:
		{
			return fmt.Sprintf("%s %s", r.SrvPriority, r.Value)
		}
	case RecordTypeMX:
		{
			return fmt.Sprintf("%s %s", r.MxPriority, r.Value)
		}
	default:
		return r.Value
	}
}
