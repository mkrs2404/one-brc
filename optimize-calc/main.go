package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"runtime/pprof"
	"sort"
)

type Calculation struct {
	Min   int
	Max   int
	Total int
	Count int
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

	// start := time.Now()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 32*1024), 1024*1024)
	var byteLine []byte
	cityWeatherMap := make(map[string]*Calculation, 1024)
	var orderedCities []string
	var cityStr string
	for scanner.Scan() {
		byteLine = scanner.Bytes()
		city, temp := parseBytes(byteLine)

		cityStr = string(city)
		calc, ok := cityWeatherMap[cityStr]
		if !ok {
			calc = &Calculation{Min: temp, Max: temp, Total: temp, Count: 1}
			cityWeatherMap[cityStr] = calc
			orderedCities = append(orderedCities, cityStr)
		} else {
			if temp < calc.Min {
				calc.Min = temp
			}
			if temp > calc.Max {
				calc.Max = temp
			}
			calc.Total += temp
			calc.Count++
		}
	}

	sort.Strings(orderedCities)
	for _, city := range orderedCities {
		calc := cityWeatherMap[city]
		avg := calc.Total / calc.Count
		fmt.Printf("%s=%.1f/%.1f/%.1f, ", city, float32(calc.Min)/10, float32(avg)/10, float32(calc.Max)/10)
	}
	// fmt.Println("Read using bufio.Scanner took :", time.Since(start).Milliseconds())
	memFile, _ := os.Create("mem.prof")
	pprof.WriteHeapProfile(memFile)
	memFile.Close()
}

func parseBytes(line []byte) ([]byte, int) {
	ind := bytes.IndexByte(line, ';')
	city := line[:ind]
	isNeg := false
	if line[ind+1] == '-' {
		ind++
		isNeg = true
	}

	splitLine := line[ind+1:]
	dp := bytes.IndexByte(splitLine, '.')

	result := int(splitLine[0] - '0')
	for i := 1; i < dp; i++ {
		result = (result << 3) + (result << 1)
		result += int(splitLine[i] - '0')
	}
	if isNeg {
		return city, -((result << 3) + (result << 1) + int(splitLine[dp+1]-'0'))
	}
	return city, (result << 3) + (result << 1) + int(splitLine[dp+1]-'0')
}
