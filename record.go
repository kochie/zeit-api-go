package zeit_api_go

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
	Id      string
	Slug    string
	Type    string
	Name    string
	Value   string
	Creator string
	Created *Time
	Updated *Time
}
