package zeit_api_go

import (
	"errors"
	"strconv"
	"time"
)

// Time is a wrapper struct to allow json unmarshal of unix timestamp
type Time struct {
	time.Time
}

// UnmarshalJSON overrides the default JSON parsing for time struct, expects the byte array to represent a unix
// timestamp with millisecond accuracy
func (t *Time) UnmarshalJSON(data []byte) error {
	unixTime, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	if unixTime < 0 {
		return errors.New("couldn't parse time")
	}
	t.Time = time.Unix(0, unixTime*1e6)
	return nil
}
