package holiday_test

import (
	"encoding/json"
	"sync"
	"testing"

	"github.com/wxb/got/holiday"
)

func TestQuery(t *testing.T) {
	wg := sync.WaitGroup{}
	years := []int{2013, 2014, 2015, 2016, 2017, 2018, 2019, 2020}
	for _, y := range years {
		wg.Add(1)

		go func(y int, wg *sync.WaitGroup) {
			defer wg.Done()

			YearHolidList, err := holiday.Query(y)
			if err != nil {
				t.Errorf("错误：%d年 %v", y, err)
			}

			for _, v := range YearHolidList {
				bytes, err := json.Marshal(v)
				if err != nil {
					t.Errorf("错误：%d年 %v", y, err)
				}

				t.Logf("%d年假日: %s", y, bytes)
			}
		}(y, &wg)
	}

	wg.Wait()
}
