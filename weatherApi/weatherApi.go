package weatherApi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/dariubs/percent"
)

type Day struct {
	Avgtemp_c    float64 `json:"avgtemp_c"`
	ChanceOfRain float64 `json:"daily_chance_of_rain"`
	Maxtemp_c    float64 `json:"maxtemp_c"`
	Mintemp_c    float64 `json:"mintemp_c"`
	Avghumidity  float64 `json:"avghumidity"`
	Maxwind_kph  float64 `json:"maxwind_kph"`
}

type ForecastDay struct {
	Date string `json:"date"`
	Day  Day    `json:"day"`
}

type Return struct {
	Forecasts []ForecastDay `json:"forecastday"`
}

type ResponseNew struct {
	Forecast Return `json:"forecast"`
}

func getData(city string, days int) ResponseNew {
	url := "http://api.weatherapi.com/v1/forecast.json?key=1444310f987149a7a31125058220812&q=" + city + "&days=" + strconv.Itoa(days)
	response, _ := http.Get(url)

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data ResponseNew
	json.Unmarshal(responseData, &data)

	return data
}

func WillItRain(city string, days int) float64 {
	if days > 3 {
		days = 3
	}

	data := getData(city, days)
	chance := 100.0
	for i := 0; i < days; i++ {
		chance -= percent.PercentFloat(data.Forecast.Forecasts[i].Day.ChanceOfRain, chance)
	}
	return 100.0 - chance
}

func RainChance(city string) []float64 {

	data := getData(city, 3)
	var chances []float64
	for i := 0; i < len(data.Forecast.Forecasts); i++ {
		chances = append(chances, float64(data.Forecast.Forecasts[i].Day.ChanceOfRain))
	}
	return chances
}

func DailyTemp(city string, days int) []float64 {
	if days > 3 {
		days = 3
	}

	data := getData(city, days)

	var returnValue []float64
	for i := 0; i < days; i++ {
		returnValue = append(returnValue, data.Forecast.Forecasts[i].Day.Avgtemp_c)
	}
	return returnValue
}

func TempArray(city string, days int) []struct {
	Day string
	Min int
	Max int
} {
	if days > 3 {
		days = 3
	}

	data := getData(city, days)

	returnData := make([]struct {
		Day string
		Min int
		Max int
	}, 0)
	for i := 0; i < days; i++ {
		returnData = append(returnData, struct {
			Day string
			Min int
			Max int
		}{data.Forecast.Forecasts[i].Date, int(data.Forecast.Forecasts[i].Day.Mintemp_c), int(data.Forecast.Forecasts[i].Day.Maxtemp_c)})
	}
	return returnData
}

func WeatherSummary(city string) struct {
	Day           string
	Min           int
	Max           int
	Humidity      int
	Wind          int
	Percipitation int
} {
	data := getData(city, 1)
	return struct {
		Day           string
		Min           int
		Max           int
		Humidity      int
		Wind          int
		Percipitation int
	}{data.Forecast.Forecasts[0].Date, int(data.Forecast.Forecasts[0].Day.Mintemp_c), int(data.Forecast.Forecasts[0].Day.Maxtemp_c),
		int(data.Forecast.Forecasts[0].Day.Avghumidity), int(data.Forecast.Forecasts[0].Day.Maxwind_kph), int(data.Forecast.Forecasts[0].Day.ChanceOfRain)}
}
