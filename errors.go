package zeit

const ErrorOrigin = "zeit API does not use `@` to represent the origin, use empty string instead"
const ErrorNilRecord = "pointer to record is nil"

type BasicError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ConflictError struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	OldId   string   `json:"oldId"`
	OldIds  []string `json:"oldIds"`
}

type GetError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Name    string `json:"name"`
}

type VerificationError struct {
	Code           string `json:"code"`
	Message        string `json:"message"`
	Name           string `json:"name"`
	NsVerification struct {
		Name                string   `json:"name"`
		Nameservers         []string `json:"nameservers"`
		IntendedNameservers []string `json:"intendedNameservers"`
	} `json:"nsVerification"`
	TxtVerification struct {
		Name               string   `json:"name"`
		Values             []string `json:"values"`
		VerificationRecord string   `json:"verificationRecord"`
	} `json:"txtVerification"`
}

type RateLimitError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Limit   struct {
		Total     int
		Remaining int
		Reset     int64
	}
}
