package openWeatherApi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/dariubs/percent"
)

type locationInfo struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

type WeatherResponse struct {
	List []OpenWeatherHour `json:"list"`
}

type OpenWeatherHour struct {
	Main struct {
		TempMin  float64 `json:"temp_min"`
		TempMax  float64 `json:"temp_max"`
		Humidity float64 `json:"humidity"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	ChanceOfRain float64 `json:"pop"`
	Date         string  `json:"dt_txt"`
}

func GetLatAndLong(city string) locationInfo {
	response, err := http.Get("http://api.openweathermap.org/geo/1.0/direct?q= " + city + "&limit=1&appid=cb79a651a59b2aff766661ae8c8f1c9e")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data []locationInfo
	json.Unmarshal(responseData, &data)

	return data[0]
}

func GetData(city string) WeatherResponse {
	loc := GetLatAndLong(city)
	url := "http://api.openweathermap.org/data/2.5/forecast?lat=" + fmt.Sprintf("%f", loc.Latitude) + "&lon=" + fmt.Sprintf("%f", loc.Longitude) + "&appid=cb79a651a59b2aff766661ae8c8f1c9e"
	response, _ := http.Get(url)

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data WeatherResponse
	json.Unmarshal(responseData, &data)

	return data
}

func ByDay(res WeatherResponse) WeatherResponse {
	var returnData WeatherResponse
	var byDay []OpenWeatherHour
	for i := 0; i < len(res.List); i++ {
		if len(byDay) == 0 {
			byDay = append(byDay, res.List[i])
		} else if byDay[0].Date[:10] == res.List[i].Date[:10] {
			byDay = append(byDay, res.List[i])
		} else {
			returnData.List = append(returnData.List, calcDay(byDay))
			byDay = make([]OpenWeatherHour, 0)
			byDay = append(byDay, res.List[i])
		}
	}
	returnData.List = append(returnData.List, calcDay(byDay))
	return returnData
}

func calcDay(hours []OpenWeatherHour) OpenWeatherHour {
	var returnData OpenWeatherHour
	returnData.Date = hours[0].Date[:10]
	for i := 0; i < len(hours); i++ {
		returnData.ChanceOfRain += hours[i].ChanceOfRain * 100
		returnData.Wind.Speed += hours[i].Wind.Speed
		returnData.Main.Humidity += hours[i].Main.Humidity
		if returnData.Main.TempMax == 0 || returnData.Main.TempMax < hours[i].Main.TempMax {
			returnData.Main.TempMax = hours[i].Main.TempMax
		}
		if returnData.Main.TempMin == 0 || returnData.Main.TempMin > hours[i].Main.TempMin {
			returnData.Main.TempMin = hours[i].Main.TempMin
		}
	}
	returnData.ChanceOfRain /= float64(len(hours))
	returnData.Wind.Speed /= float64(len(hours))
	returnData.Main.Humidity /= float64(len(hours))
	returnData.Main.TempMax -= 273.15
	returnData.Main.TempMin -= 273.15
	return returnData
}

func WillItRain(city string, days int) float64 {
	if days > 5 {
		days = 5
	}
	data := GetData(city)
	data = ByDay(data)

	chance := 100.0
	for i := 0; i < days; i++ {
		chance -= percent.PercentFloat(data.List[i].ChanceOfRain, chance)
	}
	return 100.0 - chance
}

func RainChance(city string) []float64 {
	data := GetData(city)
	data = ByDay(data)

	var chances []float64
	for i := 0; i < len(data.List); i++ {
		chances = append(chances, float64(data.List[i].ChanceOfRain))
	}
	return chances
}

func DailyTemp(city string, days int) []float64 {
	if days > 5 {
		days = 5
	}
	data := GetData(city)
	data = ByDay(data)

	var returnValue []float64
	for i := 0; i < days; i++ {
		returnValue = append(returnValue, (data.List[i].Main.TempMax+data.List[i].Main.TempMin)/2)
	}
	return returnValue
}

func TempArray(city string, days int) []struct {
	Day string
	Min int
	Max int
} {
	if days > 5 {
		days = 5
	}
	data := GetData(city)
	data = ByDay(data)

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
		}{data.List[i].Date, int(data.List[i].Main.TempMin), int(data.List[i].Main.TempMax)})
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
	data := GetData(city)
	data = ByDay(data)

	return struct {
		Day           string
		Min           int
		Max           int
		Humidity      int
		Wind          int
		Percipitation int
	}{data.List[0].Date, int(data.List[0].Main.TempMin), int(data.List[0].Main.TempMax),
		int(data.List[0].Main.Humidity), int(data.List[0].Wind.Speed), int(data.List[0].ChanceOfRain)}
}
