package zeit

const ErrorOrigin = "zeit API does not use `@` to represent the origin, use empty string instead"
const ErrorNilRecord = "pointer to record is nil"

type BasicError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e BasicError) Error() string {
	return e.Message
}

type ConflictError struct {
	BasicError
	OldId  string   `json:"oldId"`
	OldIds []string `json:"oldIds"`
}

type GetError struct {
	BasicError
	Name string `json:"name"`
}

type VerificationError struct {
	BasicError
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
	BasicError
	Limit struct {
		Total     int
		Remaining int
		Reset     int64
	}
}
