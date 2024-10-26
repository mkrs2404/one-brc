package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/pprof"
	"sort"
)

type Calculation struct {
	Min   int32
	Max   int32
	Total int64
	Count int32
}

func main() {
	f, _ := os.Create("cpu.prof")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	filePath := "../measurements.txt"
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	cityWeatherMap := make(map[string]*Calculation, 1024)
	var orderedCities []string

	scanner := bufio.NewScanner(f)
	// scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	for scanner.Scan() {
		city, temp := parseLine(scanner.Text())
		if temp == 0 {
			continue
		}

		calc, ok := cityWeatherMap[city]
		if !ok {
			calc = &Calculation{Min: temp, Max: temp, Total: int64(temp), Count: 1}
			cityWeatherMap[city] = calc
			orderedCities = append(orderedCities, city)
		} else {
			if temp < calc.Min {
				calc.Min = temp
			}
			if temp > calc.Max {
				calc.Max = temp
			}
			calc.Total += int64(temp)
			calc.Count++
		}
	}

	sort.Strings(orderedCities)
	for _, city := range orderedCities {
		calc := cityWeatherMap[city]
		avg := float32(calc.Total) / float32(calc.Count) / 10.0
		fmt.Printf("%s=%.1f/%.1f/%.1f, ", city, float32(calc.Min)/10, avg, float32(calc.Max)/10)
	}

	memFile, _ := os.Create("mem.prof")
	pprof.WriteHeapProfile(memFile)
	f.Close()
}

func parseLine(line string) (string, int32) {
	i := 0
	for i < len(line) && line[i] != ';' {
		i++
	}
	if i == len(line) {
		return "", 0
	}

	city := line[:i]
	i++ // skip semicolon

	var temp int32
	var dotSeen bool
	var decimals int

	for ; i < len(line); i++ {
		c := line[i]
		if c == '.' {
			dotSeen = true
			continue
		}
		if c < '0' || c > '9' {
			return "", 0
		}
		temp = temp*10 + int32(c-'0')
		if dotSeen {
			decimals++
			if decimals == 2 {
				break
			}
		}
	}

	if !dotSeen {
		temp *= 10
	} else if decimals == 1 {
		temp *= 10
	}

	return city, temp
}
