package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Calculation struct {
	Min   float32
	Max   float32
	Total float32
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

		temp64, err := strconv.ParseFloat(splitData[1], 32)
		if err != nil {
			continue
		}

		temp := float32(temp64)

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
		avg := calc.Total / float32(calc.Count)
		fmt.Printf("%s=%.1f/%.1f/%.1f, ", city, calc.Min, avg, calc.Max)
	}
}
