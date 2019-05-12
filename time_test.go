package zeit

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"strconv"
	"testing"
)

func TestTime_UnmarshalJSON(t *testing.T) {
	a := assert.New(t)
	unixTimes := []int64{
		0,
		123456789,
		1000000000,
		math.MaxInt32,
		math.MaxInt32 + 1,
		math.MaxInt16,
		math.MaxInt8,
		10,
		3000000,
	}
	for _, unixTime := range unixTimes {
		t.Run(fmt.Sprintf("testing time %d", unixTime), func(t *testing.T) {
			time := Time{}
			b := []byte(strconv.FormatInt(unixTime*1000, 10))

			err := time.UnmarshalJSON(b)
			a.Nil(err, "shouldn't error on correct unix timestamp parsing")
			a.Equal(unixTime, time.Unix(), "should be the same as timestamp")
		})
	}
	invalidUnixTimes := []int64{
		-1,
		10000000000000000,
		math.MaxInt64,
	}
	for _, unixTime := range invalidUnixTimes {
		t.Run(fmt.Sprintf("testing time %d", unixTime), func(t *testing.T) {
			time := Time{}
			b := []byte(strconv.FormatInt(unixTime*1000, 10))

			err := time.UnmarshalJSON(b)
			a.NotNil(err, "should error on invalid timestamp")
		})
	}
}
