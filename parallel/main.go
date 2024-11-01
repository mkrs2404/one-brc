package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"runtime/pprof"
	"sort"
	"sync"
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

	numWorkers := runtime.NumCPU()

	wg := &sync.WaitGroup{}
	byteLinesChan := make(chan []byte, 1024)
	tempChan := make(chan map[string]*Calculation, 1024)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go processData(byteLinesChan, tempChan, wg)
	}

	go func() {
		wg.Wait()
		close(tempChan)
	}()

	var readWg sync.WaitGroup
	readWg.Add(1)

	// Collect the results
	go func() {
		defer readWg.Done()
		cityWeatherMap := make(map[string]*Calculation, 512)
		orderedCities := make([]string, 0, 512)

		for cityWeather := range tempChan {
			for city, incomingCalc := range cityWeather {
				calc, ok := cityWeatherMap[city]
				if !ok {
					cityWeatherMap[city] = incomingCalc
					orderedCities = append(orderedCities, city)
				} else {
					if incomingCalc.Min < calc.Min {
						calc.Min = incomingCalc.Min
					}
					if incomingCalc.Max > calc.Max {
						calc.Max = incomingCalc.Max
					}
					calc.Total += incomingCalc.Total
					calc.Count += incomingCalc.Count
				}
			}
		}

		sort.Strings(orderedCities)
		for _, city := range orderedCities {
			calc := cityWeatherMap[city]
			avg := calc.Total / calc.Count
			fmt.Printf("%s=%.1f/%.1f/%.1f, ", city, float32(calc.Min)/10, float32(avg)/10, float32(calc.Max)/10)
		}
	}()

	bufSize := 512 * 1024
	fileBuffer := make([]byte, bufSize)
	leftover := []byte{}

	for {
		bytesRead, err := f.Read(fileBuffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		chunk := append(leftover, fileBuffer[:bytesRead]...)
		lastNewLine := bytes.LastIndexByte(chunk, '\n')
		if lastNewLine == -1 {
			leftover = chunk
			continue
		}

		// Process lines up to the last newline
		byteLines := chunk[:lastNewLine+1]
		leftover = chunk[lastNewLine+1:]

		byteLinesChan <- byteLines
	}

	close(byteLinesChan)

	wg.Wait()
	readWg.Wait()
	memFile, _ := os.Create("mem.prof")
	pprof.WriteHeapProfile(memFile)
	memFile.Close()
}

func processData(byteChan <-chan []byte, tempChan chan<- map[string]*Calculation, wg *sync.WaitGroup) {
	defer wg.Done()

	for byteLines := range byteChan {
		orderedCities := make([]string, 0, 1024)
		cityWeatherMap := make(map[string]*Calculation, 1024)
		var byteLine []byte
		var cityStr string

		// Process each line
		start := 0
		for i := 0; i < len(byteLines); i++ {
			if byteLines[i] == 10 { // newline found
				byteLine = byteLines[start:i]
				start = i + 1

				index, temp := parseBytes(byteLine)

				cityStr = string(byteLine[:index])
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
		}
		tempChan <- cityWeatherMap
	}
}

func parseBytes(line []byte) (int, int) {
	var ind int
	for i, v := range line {
		if v == ';' {
			ind = i
			break
		}
	}

	cityInd := ind
	isNeg := false
	if line[ind+1] == '-' {
		ind++
		isNeg = true
	}

	temp := line[ind+1:]
	var dp int
	for i, v := range temp {
		if v == '.' {
			dp = i
			break
		}
	}

	result := int(temp[0] - '0')
	for i := 1; i < dp; i++ {
		result = (result << 3) + (result << 1)
		result += int(temp[i] - '0')
	}

	finalTemp := ((result << 3) + (result << 1) + int(temp[dp+1]-'0'))
	if isNeg {
		finalTemp = -finalTemp
	}
	return cityInd, finalTemp
}