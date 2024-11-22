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
	"unsafe"
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

	// Larger buffer sizes for channels to reduce contention
	byteLinesChan := make(chan []byte, numWorkers*2)
	collectorChan := make(chan map[string]*Calculation, numWorkers*2)

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	// Start workers
	for i := 0; i < numWorkers; i++ {
		go processData(byteLinesChan, collectorChan, &wg)
	}

	// Collector goroutine
	var readWg sync.WaitGroup
	readWg.Add(1)

	go func() {
		defer readWg.Done()
		results := make(map[string]*Calculation, 512)
		orderedCities := make([]string, 0, 512)
		citiesMap := make(map[string]struct{}, 512)

		for cityWeather := range collectorChan {
			for city, incomingCalc := range cityWeather {
				if calc, ok := results[city]; !ok {
					results[city] = incomingCalc
					if _, exists := citiesMap[city]; !exists {
						orderedCities = append(orderedCities, city)
						citiesMap[city] = struct{}{}
					}
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

		// Output results
		sort.Strings(orderedCities)
		for _, city := range orderedCities {
			calc := results[city]
			avg := float32(calc.Total) / float32(calc.Count)
			fmt.Printf("%s=%.1f/%.1f/%.1f, ", city, float32(calc.Min)/10, avg/10, float32(calc.Max)/10)
		}
	}()

	const bufSize = 1024 * 1024
	buf := make([]byte, bufSize)
	leftover := make([]byte, 0, 1024)

	for {
		n, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		chunk := append(leftover, buf[:n]...)
		lastNL := bytes.LastIndexByte(chunk, '\n')
		if lastNL == -1 {
			leftover = append(leftover[:0], chunk...)
			continue
		}

		byteLinesChan <- chunk[:lastNL+1]
		leftover = append(leftover[:0], chunk[lastNL+1:]...)
	}

	if len(leftover) > 0 {
		byteLinesChan <- leftover
	}

	close(byteLinesChan)
	wg.Wait()
	close(collectorChan)
	readWg.Wait()
	memFile, _ := os.Create("mem.prof")
	pprof.WriteHeapProfile(memFile)
	memFile.Close()
}

func processData(byteChan <-chan []byte, tempChan chan<- map[string]*Calculation, wg *sync.WaitGroup) {

	defer wg.Done()

	for chunk := range byteChan {
		cityWeatherMap := make(map[string]*Calculation, 128)

		start := 0
		for i := 0; i < len(chunk); i++ {
			if chunk[i] == '\n' {
				line := chunk[start:i]
				cityEnd, temp := parseBytes(line)

				// Zero-allocation string conversion
				cityStr := unsafe.String(unsafe.SliceData(line), cityEnd)

				if calc, ok := cityWeatherMap[cityStr]; !ok {
					cityWeatherMap[cityStr] = &Calculation{
						Min:   temp,
						Max:   temp,
						Total: temp,
						Count: 1,
					}
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
				start = i + 1
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
