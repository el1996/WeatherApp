package timeAndDate

import (
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/dariubs/percent"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
)

type Table struct {
	Date      []string
	TempMin   []int
	TempMax   []int
	FeelsLike []int
	Wind      []int
	Humidity  []int
	Precip    []int
	Amount    []int
}

var tab Table

func Init(city string) {
	geziyor.NewGeziyor(&geziyor.Options{
		StartRequestsFunc: func(g *geziyor.Geziyor) {
			g.GetRendered("https://www.timeanddate.com/weather/israel/"+city+"/ext", g.Opt.ParseFunc)
		},
		ParseFunc: quotesParse,
		//Exporters: []export.Exporter{&export.JSON{}},
	}).Start()

}

func quotesParse(g *geziyor.Geziyor, r *client.Response) {
	r.HTMLDoc.Find("#wt-ext").Each(func(i int, s *goquery.Selection) {
		s.Find("tr").Each(func(indexRow int, sRow *goquery.Selection) {
			if indexRow == 0 || indexRow == 1 || indexRow == 17 {
				return
			}
			day := sRow.Find("th").First()
			tab.Date = append(tab.Date, stringToDate(day.Text()))
			sRow.Find("td").Each(func(indexCol int, sCol *goquery.Selection) {
				if indexCol == 0 || indexCol == 2 || indexCol == 5 {
					return
				}
				if indexCol == 1 {
					firstSplit := strings.Split(sCol.Text(), " ")
					value, err := strconv.Atoi(firstSplit[0])
					if err != nil {
						return
					}
					tab.TempMax = append(tab.TempMax, value)
					secondSplit := strings.Split(firstSplit[2], "\u00a0")
					value, err = strconv.Atoi(strings.Trim(secondSplit[0], " "))
					if err != nil {
						return
					}
					tab.TempMin = append(tab.TempMin, value)
				}
				if indexCol == 3 {
					firstSplit := strings.Split(sCol.Text(), "\u00a0")
					value, _ := strconv.Atoi(strings.Trim(firstSplit[0], " "))
					tab.FeelsLike = append(tab.FeelsLike, value)
				}
				if indexCol == 4 {
					firstSplit := strings.Split(sCol.Text(), " ")
					value, _ := strconv.Atoi(strings.Trim(firstSplit[0], " "))
					tab.Wind = append(tab.Wind, value)
				}
				if indexCol == 6 {
					value, _ := strconv.Atoi(strings.Trim(sCol.Text(), "%"))
					tab.Humidity = append(tab.Humidity, value)
				}
				if indexCol == 7 {
					value, _ := strconv.Atoi(strings.Trim(sCol.Text(), "%"))
					tab.Precip = append(tab.Precip, value)
				}
				if indexCol == 8 {
					if sCol.Text() == "-" {
						tab.Amount = append(tab.Amount, 0)
					} else {
						firstSplit := strings.Split(sCol.Text(), ".")
						value, _ := strconv.Atoi(strings.Trim(firstSplit[0], " "))
						tab.Amount = append(tab.Amount, value)
					}
				}
			})
		})
	})
}

func stringToDate(date string) string {
	var day int
	day, _ = strconv.Atoi(date[7:9])
	var month time.Month
	switch date[3:6] {
	case "Jan":
		month = time.January
	case "Feb":
		month = time.February
	case "Mar":
		month = time.March
	case "Apr":
		month = time.April
	case "May":
		month = time.May
	case "Jun":
		month = time.June
	case "Jul":
		month = time.July
	case "Aug":
		month = time.August
	case "Sep":
		month = time.September
	case "Oct":
		month = time.October
	case "Nov":
		month = time.November
	case "Dec":
		month = time.December
	}
	return time.Date(time.Now().Year(), month, day, 0, 0, 0, 0, time.UTC).Local().String()[0:10]
}

func WillItRain(days int) float64 {
	if days > 14 {
		days = 14
	}
	chance := 100.0
	for i := 0; i < days; i++ {
		chance -= percent.PercentFloat(float64(tab.Precip[i]), chance)
	}
	return 100.0 - chance
}

func RainChance() []float64 {
	var chances []float64
	for i := 0; i < len(tab.Precip); i++ {
		chances = append(chances, float64(tab.Precip[i]))
	}
	return chances
}

func DailyTemp(days int) []float64 {
	if days > 14 {
		days = 14
	}

	var returnValue []float64
	for i := 0; i < days; i++ {
		returnValue = append(returnValue, float64((tab.TempMax[i]+tab.TempMin[i]))/2)
	}
	return returnValue
}

func TempArray(days int) []struct {
	Day string
	Min int
	Max int
} {
	if days > 14 {
		days = 14
	}
	data := make([]struct {
		Day string
		Min int
		Max int
	}, 0)
	for i := 0; i < days; i++ {
		data = append(data, struct {
			Day string
			Min int
			Max int
		}{tab.Date[i], tab.TempMin[i], tab.TempMax[i]})
	}
	return data
}

func WeatherSummary() struct {
	Day           string
	Min           int
	Max           int
	Humidity      int
	Wind          int
	Percipitation int
} {
	return struct {
		Day           string
		Min           int
		Max           int
		Humidity      int
		Wind          int
		Percipitation int
	}{tab.Date[0], tab.TempMin[0], tab.TempMax[0], tab.Humidity[0], tab.Wind[0], tab.Precip[0]}
}
