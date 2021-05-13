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

var year string
var years []string

var c *colly.Collector

func init() {
	c = colly.NewCollector()

	c.OnRequest(onRequest)
	c.OnResponse(onResponse)
	c.OnError(onError)
}

type H struct {
	year        string
	reqCallback colly.RequestCallback
	rspCallback colly.ResponseCallback
}

func onRequest(r *colly.Request) {
	fmt.Println("Visiting", r.URL)
}

func onResponse(r *colly.Response) {
}

func onError(r *colly.Response, err error) {

}

func Query(y string) ([]Holiday, error) {
	year = y
	api := "https://tool.lu/holiday/index.html?y=" + year
	YearHolidList := []Holiday{}

	c.OnHTML("div.inner > table.tbl", func(e *colly.HTMLElement) {
		e.DOM.Find("tr").Each(func(i int, e *goquery.Selection) {
			if i != 0 {

				wd := []rune(e.Find("td").Eq(3).Text())
				wdi, _ := strconv.Atoi(string(wd[:len(wd)-1]))

				se := strings.Split(e.Find("td").Eq(1).Text(), "~")

				timeLayout := "2006年01月02日"
				loc, _ := time.LoadLocation("Local")
				sdate, _ := time.ParseInLocation(timeLayout, year+"年"+se[0], loc)
				edate, _ := time.ParseInLocation(timeLayout, year+"年"+se[1], loc)

				YearHolidList = append(YearHolidList, Holiday{
					Name:    e.Find("td").Eq(0).Text(),
					Sdate:   sdate,
					Edate:   edate,
					WorkDay: e.Find("td").Eq(2).Text(),
					DayNum:  wdi,
					DayUnit: string(wd[len(wd)-1:]),
				})
			}

		})

		fmt.Println(len(YearHolidList))
	})

	err := c.Visit(api)

	return YearHolidList, err
}
