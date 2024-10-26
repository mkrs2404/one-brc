package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/pprof"
	"time"
)

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

	// Read the file using bufio.Scanner
	start := time.Now()
	// readScannerStr(f)
	readScannerBytes(f)
	fmt.Println("Read using bufio.Scanner took :", time.Since(start).Milliseconds())
	memFile, _ := os.Create("mem.prof")
	pprof.WriteHeapProfile(memFile)
	memFile.Close()
}

// Read the file using bufio.Scanner
func readScannerStr(f *os.File) {
	scanner := bufio.NewScanner(f)
	// scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	var line string

	for scanner.Scan() {
		line = scanner.Text()
	}
	_ = line
}

// Read the file using bufio.Scanner
func readScannerBytes(f *os.File) {
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 32*1024), 1024*1024)
	var byteLine []byte

	for scanner.Scan() {
		byteLine = scanner.Bytes()
	}
	_ = byteLine
}