package holiday_test

import (
	"fmt"
	"testing"

	"github.com/wxb/got/holiday"
)

func TestQuery(t *testing.T) {
	YearHolidList, err := holiday.Query("2021")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(YearHolidList)
}
