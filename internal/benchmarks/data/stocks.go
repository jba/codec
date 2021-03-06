// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package data

import (
	"encoding/gob"
	"math"
	"math/rand"
	"os"
	"time"
)

// Stock price data.
//
// Idea from https://towardsdatascience.com/how-to-store-financial-market-data-for-backtesting-84b95fc016fc.
//
// This is not the most efficient representation, but it does demonstrate
// numbers-heavy data.

var Stocks = gobBenchmarkData("stocks", func() interface{} { return new([]*StockData) })

type StockData struct {
	Symbol    string
	Intervals []Interval
}

func (s1 *StockData) Equal(s2 *StockData) bool {
	if s1.Symbol != s2.Symbol {
		return false
	}
	if len(s1.Intervals) != len(s2.Intervals) {
		return false
	}
	for i, v := range s1.Intervals {
		if !v.Equal(s2.Intervals[i]) {
			return false
		}
	}
	return true
}

type Interval struct {
	Start, End                     time.Time
	Open, Close, Low, High, Volume float64
}

func (i Interval) Equal(j Interval) bool {
	return i.Start.Equal(j.Start) &&
		i.End.Equal(j.End) &&
		i.Open == j.Open &&
		i.Close == j.Close &&
		i.Low == j.Low &&
		i.High == j.High &&
		i.Volume == j.Volume
}

func generateStockDataToFile(filename string) error {
	sds := generateStockData(200, 365*20)
	return WriteNewFile(filename, func(f *os.File) error {
		return gob.NewEncoder(f).Encode(sds)
	})
}

// Generate random stock data. The stock prices and volumes make no sense; this
// is not a simulation. We just need data of the right form.
func generateStockData(nStocks, nIntervals int) []*StockData {
	var sds []*StockData
	for i := 0; i < nStocks; i++ {
		sds = append(sds, generateStockData1(nIntervals))
	}
	return sds
}

func generateStockData1(n int) *StockData {
	var bytes [4]byte
	for i := 0; i < len(bytes); i++ {
		bytes[i] = byte('A' + rand.Intn(26))
	}
	start := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	return &StockData{
		Symbol:    string(bytes[:]),
		Intervals: generateIntervals(start, n),
	}
}

func generateIntervals(start time.Time, n int) []Interval {
	ivs := make([]Interval, n)
	y, m, d := start.Date()
	date := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

	pennies := func(x float64) float64 { return math.Round(100*x) / 100 }

	for i := 0; i < n; i++ {
		low := 1000 * rand.Float64()        // between $0 and $1000
		high := low + 1 + 49*rand.Float64() // low + ($1 to $50)

		ivs[i] = Interval{
			Start:  date.Add(9 * time.Hour),  // 9 AM
			End:    date.Add(17 * time.Hour), // 5 PM
			Low:    pennies(low),
			High:   pennies(high),
			Open:   pennies(low + (high-low)*rand.Float64()),
			Close:  pennies(low + (high-low)*rand.Float64()),
			Volume: rand.Float64() * 100,
		}
		date = date.AddDate(0, 0, 1)
	}
	return ivs
}
