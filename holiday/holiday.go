package holiday

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type Holiday struct {
	Name    string    `json:"name"`       // 节日名称
	Sdate   time.Time `json:"start_date"` // 节日开始日期
	Edate   time.Time `json:"end_date"`   // 节日结束日期
	WorkDay string    `json:"work_day"`   // 调班工作日
	DayNum  int       `json:"day_num"`    // 放假天数
	DayUnit string    `json:"day_unit"`   // 放假天数单位（默认：天）
}

const api = "https://tool.lu/holiday/index.html"

var c *colly.Collector

var layout string = "2006年01月02日"
var beginY time.Time
var endY time.Time

func init() {
	c = colly.NewCollector()

	beginY = time.Date(2013, time.January, 1, 0, 0, 0, 0, time.Local)
	endY = time.Date(time.Now().Year()-1, time.January, 1, 0, 0, 0, 0, time.Local)
}

func Query(year int) (YearHolidList []Holiday, err error) {
	// 参数检查
	y := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local)
	if y.Before(beginY) || y.After(endY) {
		err = fmt.Errorf("the year param must in[%d~%d]", beginY.Year(), endY.Year())
		return
	}

	c.OnHTML("div.inner > table.tbl", func(e *colly.HTMLElement) {
		e.DOM.Find("tr").Each(func(i int, e *goquery.Selection) {
			if i != 0 {
				wd := []rune(e.Find("td").Eq(3).Text())
				wdi, err := strconv.Atoi(string(wd[:len(wd)-1]))
				if err == nil {
					se := strings.Split(e.Find("td").Eq(1).Text(), "~")
					sdate, _ := time.Parse(layout, strconv.Itoa(year)+"年"+se[0])
					edate, _ := time.Parse(layout, strconv.Itoa(year)+"年"+se[1])

					YearHolidList = append(YearHolidList, Holiday{
						Name:    e.Find("td").Eq(0).Text(),
						Sdate:   sdate,
						Edate:   edate,
						WorkDay: e.Find("td").Eq(2).Text(),
						DayNum:  wdi,
						DayUnit: string(wd[len(wd)-1:]),
					})
				}

			}

		})
	})

	err = c.Visit(fmt.Sprintf("%s?y=%d", api, year))
	return
}
