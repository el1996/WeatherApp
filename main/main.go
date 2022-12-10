package main

import (
	"example.com/timeAndDate"

	"example.com/weatherApi"

	"fmt"
	"time"

	"example.com/openWeatherApi"
)

func main() {
	timeAndDate.Init("haifa")

	fmt.Println(WillItRain("haifa", 2))
	fmt.Println(WillItRain("haifa", 5))
	fmt.Println(WillItRain("haifa", 14))

	fmt.Println(NextRainDay("haifa", 2))
	fmt.Println(NextRainDay("haifa", 5))
	fmt.Println(NextRainDay("haifa", 14))

	fmt.Println(AverageTemp("haifa", 2))
	fmt.Println(AverageTemp("haifa", 5))
	fmt.Println(AverageTemp("haifa", 14))

	fmt.Println(TempArray("haifa", 2))
	fmt.Println(TempArray("haifa", 5))
	fmt.Println(TempArray("haifa", 14))

	fmt.Println(WeatherSummary("haifa"))
}

func WillItRain(city string, days int) float64 {
	if days > 5 {
		return timeAndDate.WillItRain(days)
	} else if days > 3 {
		return (timeAndDate.WillItRain(days) + openWeatherApi.WillItRain(city, days)) / 2
	} else if days > 0 {
		return (timeAndDate.WillItRain(days) + openWeatherApi.WillItRain(city, days) + weatherApi.WillItRain(city, days)) / 3
	} else {
		return 0
	}
}

func NextRainDay(city string, days int) time.Time {
	if days > 14 {
		days = 14
	}
	timeArr := timeAndDate.RainChance()
	openArr := openWeatherApi.RainChance(city)
	weatherArr := weatherApi.RainChance(city)

	for i := 0; i < days; i++ {
		if i > 4 {
			if timeArr[i] > 50 {
				return time.Now().Add(time.Hour * time.Duration(24*i))
			}
		} else if i > 2 {
			if (timeArr[i]+openArr[i])/2 > 50 {
				return time.Now().Add(time.Hour * time.Duration(24*i))
			}
		} else {
			if (timeArr[i]+openArr[i]+weatherArr[i])/3 > 50 {
				return time.Now().Add(time.Hour * time.Duration(24*i))
			}
		}
	}
	return time.Now().Add(time.Hour * time.Duration(24*-1))
}

func AverageTemp(city string, days int) float64 {
	if days > 14 {
		days = 14
	}
	timeArr := timeAndDate.DailyTemp(days)
	openArr := openWeatherApi.DailyTemp(city, days)
	weatherArr := weatherApi.DailyTemp(city, days)

	sum := 0.0
	for i := 0; i < days; i++ {
		if i > 4 {
			sum += timeArr[i]
		} else if i > 2 {
			sum += (timeArr[i] + openArr[i]) / 2
		} else {
			sum += (timeArr[i] + openArr[i] + weatherArr[i]) / 3
		}
	}
	return sum / float64(days)
}

func TempArray(city string, days int) []struct {
	Day string
	Min int
	Max int
} {
	returnValue := make([]struct {
		Day string
		Min int
		Max int
	}, 0)

	if days > 14 {
		days = 14
	}
	timeArr := timeAndDate.TempArray(days)
	openArr := openWeatherApi.TempArray(city, days)
	weatherArr := weatherApi.TempArray(city, days)

	for i := 0; i < days; i++ {
		if i > 4 {
			returnValue = append(returnValue, struct {
				Day string
				Min int
				Max int
			}{timeArr[i].Day, timeArr[i].Min, timeArr[i].Max})
		} else if i > 2 {
			returnValue = append(returnValue, struct {
				Day string
				Min int
				Max int
			}{timeArr[i].Day, (timeArr[i].Min + openArr[i].Min) / 2, (timeArr[i].Max + openArr[i].Max) / 2})
		} else {
			returnValue = append(returnValue, struct {
				Day string
				Min int
				Max int
			}{timeArr[i].Day, (timeArr[i].Min + openArr[i].Min + weatherArr[i].Min) / 3, (timeArr[i].Max + openArr[i].Max + weatherArr[i].Max) / 3})
		}
	}
	return returnValue
}

func WeatherSummary(city string) struct {
	Day           string
	Min           int
	Max           int
	Humidity      int
	Wind          int
	Percipitation int
} {
	timeAndDate := timeAndDate.WeatherSummary()
	open := openWeatherApi.WeatherSummary(city)
	weather := weatherApi.WeatherSummary(city)
	dateNow := time.Now().Local().String()[0:10]
	amount := 0
	var returnValue struct {
		Day           string
		Min           int
		Max           int
		Humidity      int
		Wind          int
		Percipitation int
	}

	if timeAndDate.Day == dateNow {
		amount++
		returnValue.Day = timeAndDate.Day
		returnValue.Min += timeAndDate.Min
		returnValue.Max += timeAndDate.Max
		returnValue.Humidity += timeAndDate.Humidity
		returnValue.Wind += timeAndDate.Wind
		returnValue.Percipitation += timeAndDate.Percipitation
	}
	if open.Day == dateNow {
		amount++
		returnValue.Day = open.Day
		returnValue.Min += open.Min
		returnValue.Max += open.Max
		returnValue.Humidity += open.Humidity
		returnValue.Wind += open.Wind
		returnValue.Percipitation += open.Percipitation
	}
	if weather.Day == dateNow {
		amount++
		returnValue.Day = weather.Day
		returnValue.Min += weather.Min
		returnValue.Max += weather.Max
		returnValue.Humidity += weather.Humidity
		returnValue.Wind += weather.Wind
		returnValue.Percipitation += weather.Percipitation
	}

	if amount == 0 {
		return struct {
			Day           string
			Min           int
			Max           int
			Humidity      int
			Wind          int
			Percipitation int
		}{}
	}

	returnValue.Min /= amount
	returnValue.Max /= amount
	returnValue.Humidity /= amount
	returnValue.Wind /= amount
	returnValue.Percipitation /= amount

	return returnValue
}
