package main

import (
	"encoding/gob"
	"math/rand"
	"os"
	"time"
)

// Synthetic, integer-heavy benchmark.

var scores = gobBenchmarkData("scores", func() interface{} { return new([]Score) })

type Score struct {
	GameID   int
	PlayerID int
	Scores   []int
}

func generateScoreDataToFile(filename string) error {
	sd := generateScoreData(10_000, 100)
	return writeNewFile(filename, func(f *os.File) error {
		return gob.NewEncoder(f).Encode(sd)
	})
}

func generateScoreData(nRecords, nScoresPerRecord int) []Score {
	z := rand.NewZipf(rand.New(rand.NewSource(time.Now().UnixNano())), 1.2, 1, 1_000_000)
	r := make([]Score, nRecords)
	for i := 0; i < len(r); i++ {
		r[i] = generateScore(nScoresPerRecord, z)
	}
	return r
}

func generateScore(n int, z *rand.Zipf) Score {
	s := make([]int, n)
	for i := 0; i < n; i++ {
		s[i] = int(z.Uint64())
	}
	return Score{
		GameID:   rand.Intn(1000),
		PlayerID: rand.Intn(10_000),
		Scores:   s,
	}
}
