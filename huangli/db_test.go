package huangli

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContain(t *testing.T) {
	holidayList := HolidayList([]*Holiday{
		{
			Date:   "2020-12-1",
			Status: 0,
		},
		{
			Date:   "2020-12-2",
			Status: 0,
		},
		{
			Date:   "2020-11-1",
			Status: 0,
		},
	})
	tests := []struct {
		InputYear  int
		InputMouth int
		Expect     bool
	}{
		{
			InputYear:  2020,
			InputMouth: 11,
			Expect:     true,
		},
		{
			InputYear:  2020,
			InputMouth: 12,
			Expect:     true,
		},
		{
			InputYear:  2020,
			InputMouth: 9,
			Expect:     false,
		},
		{
			InputYear:  2019,
			InputMouth: 12,
			Expect:     false,
		},
	}
	for _, data := range tests {
		assert.Equal(t, data.Expect, holidayList.Contain(data.InputYear, data.InputMouth))
	}
}
