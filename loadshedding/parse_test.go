package loadshedding

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_parseDate(t *testing.T) {
	assert := assert.New(t)

	tz, err := time.LoadLocation("Africa/Johannesburg")
	assert.NoError(err)

	actual, err := parseDate("Mon, 07 Sep", tz)
	assert.NoError(err)
	assert.Equal(actual, time.Date(0, time.September, 7, 0, 0, 0, 0, tz))
}

func Test_parseTime(t *testing.T) {
	assert := assert.New(t)

	tz, err := time.LoadLocation("Africa/Johannesburg")
	assert.NoError(err)

	actual, err := parseTime("02:30", tz)
	assert.NoError(err)
	assert.Equal(actual, time.Date(0, time.January, 1, 2, 30, 0, 0, tz))
}

func Test_parseHTMLTime(t *testing.T) {
	assert := assert.New(t)

	tz, err := time.LoadLocation("Africa/Johannesburg")
	assert.NoError(err)

	actualStart, actualEnd, err := parseHTMLTime(2020, "Mon, 07 Sep", "00:00 - 02:30", tz)
	assert.NoError(err)
	assert.Equal(*actualStart, time.Date(2020, time.September, 7, 0, 0, 0, 0, tz))
	assert.Equal(*actualEnd, time.Date(2020, time.September, 7, 2, 30, 0, 0, tz))
}
