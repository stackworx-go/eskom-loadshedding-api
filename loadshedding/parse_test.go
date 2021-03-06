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

	slots, err := parseHTMLTime(2020, "Mon, 07 Sep", "00:00 - 02:30", tz, 0)

	assert.NoError(err)
	assert.Equal(1, len(slots.Slots))
	assert.Equal(slots.Date, time.Date(2020, time.September, 7, 0, 0, 0, 0, tz))
	assert.Equal(slots.Slots[0].Start, time.Date(2020, time.September, 7, 0, 0, 0, 0, tz))
	assert.Equal(slots.Slots[0].Duration, mustParseDuration("2h30m"))
}

func mustParseDuration(val string) time.Duration {
	d, err := time.ParseDuration(val)

	if err != nil {
		panic(err)
	}

	return d
}

func Test_parseHTMLTime_doubleSlot(t *testing.T) {
	assert := assert.New(t)

	tz, err := time.LoadLocation("Africa/Johannesburg")
	assert.NoError(err)

	slots, err := parseHTMLTime(2020, "Mon, 07 Sep", "04:00 - 08:3020:00 - 00:30", tz, 0)

	assert.NoError(err)
	assert.Equal(2, len(slots.Slots))
	assert.Equal(slots.Date, time.Date(2020, time.September, 7, 0, 0, 0, 0, tz))
	assert.Equal(slots.Slots[0].Start, time.Date(2020, time.September, 7, 4, 0, 0, 0, tz))
	assert.Equal(slots.Slots[0].Duration, mustParseDuration("4h30m"))
	assert.Equal(slots.Slots[1].Start, time.Date(2020, time.September, 7, 20, 0, 0, 0, tz))
	assert.Equal(slots.Slots[1].Duration, mustParseDuration("4h30m"))
}
