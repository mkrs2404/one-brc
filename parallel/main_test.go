package main

import (
	"bytes"
	"fmt"
	"testing"

	"math/rand"
)

func Int() int {
	result := 2
	finalTemp := ((result << 3) + (result << 1) + 1)
	a := finalTemp / result
	return a
}

func Int32() int32 {
	result := 2
	finalTemp := int32(((result << 3) + (result << 1) + 1))
	a := finalTemp / int32(result)
	return a
}

func Int64() int64 {
	result := 2
	finalTemp := int64(((result << 3) + (result << 1) + 1))
	a := finalTemp / int64(result)
	return a
}

func Float64Add() float64 {
	result := 2.0
	finalTemp := result + 1.0
	return finalTemp
}

func IntAdd() int {
	result := 2
	finalTemp := result + 1
	return finalTemp
}

func Float64Div() float64 {
	result := 3567.0
	finalTemp := 233.0
	return result / finalTemp
}

func IntDiv() int {
	result := 3567
	finalTemp := 233
	return result / finalTemp
}

func MultBit() int {
	result := 2
	finalTemp := ((result << 3) + (result << 1))
	return finalTemp
}

func Mult() int {
	result := 2
	finalTemp := result * 10
	return finalTemp
}

func BenchmarkMult(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Mult()
	}
}

func BenchmarkMultBit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MultBit()
	}
}

func BenchmarkFloat64Add(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Float64Add()
	}
}

func BenchmarkIntAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		a := IntAdd()
		_ = a
	}
}

func BenchmarkFloat64Div(b *testing.B) {
	for i := 0; i < b.N; i++ {
		a := Float64Div()
		_ = a
	}
}

func BenchmarkIntDiv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = IntDiv()
	}
}

func BenchmarkInt(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Int()
	}
}

func BenchmarkInt32(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Int32()
	}
}

func BenchmarkInt64(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Int64()
	}
}

// manualIndexByte implements a simple loop to find byte
func manualIndexByte(b []byte, c byte) int {
	for i := 0; i < len(b); i++ {
		if b[i] == c {
			return i
		}
	}
	return -1
}

// generateTestData creates a byte slice with the target byte at position pos
func generateTestData(size, pos int, target byte) []byte {
	data := make([]byte, size)
	// Fill with random ASCII letters
	for i := range data {
		data[i] = byte(rand.Intn(26) + 'a')
	}
	if pos < size {
		data[pos] = target
	}
	return data
}

func BenchmarkBytesSearch(b *testing.B) {
	sizes := []int{8, 16, 32, 64, 128, 256, 512, 1024}
	positions := []float64{0.25, 0.5, 0.75} // Search positions as percentage of size

	for _, size := range sizes {
		for _, pos := range positions {
			position := int(float64(size) * pos)
			data := generateTestData(size, position, ';')

			b.Run(fmt.Sprintf("stdlib/size=%d/pos=%d", size, position), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					bytes.IndexByte(data, ';')
				}
			})

			b.Run(fmt.Sprintf("manual/size=%d/pos=%d", size, position), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					manualIndexByte(data, ';')
				}
			})
		}
	}
}

func BenchmarkWorstCase(b *testing.B) {
	sizes := []int{8, 64, 512}

	for _, size := range sizes {
		data := generateTestData(size, -1, ';') // No target byte present

		b.Run("stdlib/size="+string(rune(size)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bytes.IndexByte(data, ';')
			}
		})

		b.Run("manual/size="+string(rune(size)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				manualIndexByte(data, ';')
			}
		})
	}
}

var Result int

func BenchmarkPreventOptimization(b *testing.B) {
	data := generateTestData(64, 32, ';')
	var r int

	b.Run("stdlib", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r = bytes.IndexByte(data, ';')
		}
		Result = r
	})

	b.Run("manual", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r = manualIndexByte(data, ';')
		}
		Result = r
	})
}
