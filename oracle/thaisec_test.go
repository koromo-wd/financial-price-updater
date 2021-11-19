package oracle

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetQueryNavDate(t *testing.T) {
	now = time.Date(2021, time.October, 31, 0, 0, 0, 0, time.UTC)

	assert.Equal(t, "2021-10-31", getQueryNavDate(0))
	assert.Equal(t, "2021-10-30", getQueryNavDate(1))
	assert.Equal(t, "2021-09-30", getQueryNavDate(31))
}

func TestGetTimeLoc(t *testing.T) {
	bkk, _ := time.LoadLocation("Asia/Bangkok")

	result, err := getTimeLoc("Asia/Bangkok")
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, bkk, result)

	utc, _ := time.LoadLocation("UTC")

	result, err = getTimeLoc("")
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, utc, result)
}
