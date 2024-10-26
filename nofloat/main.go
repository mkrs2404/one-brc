package main

import (
	"bufio"
	"fmt"
	_ "net/http/pprof"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Calculation struct {
	Min   int64
	Max   int64
	Total int64
	Count int64
}

func main() {
	filePath := "../measurements.txt"

	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	cityWeatherMap := make(map[string]Calculation)
	orderedCities := make([]string, 0, 1024)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		splitData := strings.Split(line, ";")
		if len(splitData) != 2 {
			continue
		}

		modifiedTempStr := strings.Replace(splitData[1], ".", "", 1)
		temp, err := strconv.ParseInt(modifiedTempStr, 10, 32)
		if err != nil {
			continue
		}

		calc, ok := cityWeatherMap[splitData[0]]
		if !ok {
			cityWeatherMap[splitData[0]] = Calculation{
				Min:   temp,
				Max:   temp,
				Total: temp,
				Count: 1,
			}
			orderedCities = append(orderedCities, splitData[0])
		} else {
			if temp < calc.Min {
				calc.Min = temp
			}
			if temp > calc.Max {
				calc.Max = temp
			}
			calc.Total += temp
			calc.Count++
			cityWeatherMap[splitData[0]] = calc
		}
	}
	sort.Strings(orderedCities)
	for _, city := range orderedCities {
		calc := cityWeatherMap[city]
		avg := (float32(calc.Total) / 10) / float32(calc.Count)
		fmt.Printf("%s=%.1f/%.1f/%.1f, ", city, float32(calc.Min)/10, avg, float32(calc.Max)/10)
	}
}
