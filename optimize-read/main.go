package main

import (
	"bufio"
	"fmt"
	"io"
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
	// readScannerBytes(f)
	readByFileReader(f)
	fmt.Println("Read using file reader took :", time.Since(start).Milliseconds())
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

func readByFileReader(f *os.File) {
	// pos, neg := 0, 0
	buffer := make([]byte, 2048*2048)
	leftover := []byte{}
	for {
		bytesRead, err := f.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		// Combine leftover with current buffer
		chunk := append(leftover, buffer[:bytesRead]...)

		lastNewLine := -1
		for i := len(chunk) - 1; i >= 0; i-- {
			if chunk[i] == 10 { //10 is the ASCII value of newline
				lastNewLine = i
				break
			}
		}

		var processUntil int
		// No newline found, entire chunk becomes leftover
		if lastNewLine == -1 {
			leftover = chunk
			continue
		} else {
			processUntil = lastNewLine + 1
			leftover = chunk[processUntil:]
		}

		byteLines := chunk[:processUntil]

		bufIndex := 0
		for {
			// Extracting a single line from the multiple lines inside the byteLines slice
			var start, end int
			for i := bufIndex; i < len(byteLines); i++ {
				if byteLines[i] == 10 {
					start, end = bufIndex, i
					bufIndex = i + 1
					break
				}
			}

			if bufIndex == len(byteLines) {
				break
			}
			_ = byteLines[start:end]
		}
	}
}
