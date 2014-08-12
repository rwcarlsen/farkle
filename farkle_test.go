package farkle

import (
	"math/rand"
	"testing"
	"time"
)

type roll struct {
	Dice
	Score int
	Nrem  int
}

var rolls = []roll{
	// double triples
	roll{Dice{1: 6, 2: 0, 3: 0, 4: 0, 5: 0, 6: 0}, 2000, 0},
	roll{Dice{1: 0, 2: 6, 3: 0, 4: 0, 5: 0, 6: 0}, 400, 0},
	roll{Dice{1: 0, 2: 0, 3: 6, 4: 0, 5: 0, 6: 0}, 600, 0},
	roll{Dice{1: 0, 2: 0, 3: 0, 4: 6, 5: 0, 6: 0}, 800, 0},
	roll{Dice{1: 0, 2: 0, 3: 0, 4: 0, 5: 6, 6: 0}, 1000, 0},
	roll{Dice{1: 0, 2: 0, 3: 0, 4: 0, 5: 0, 6: 6}, 1200, 0},
	roll{Dice{1: 0, 2: 0, 3: 0, 4: 0, 5: 0, 6: 6}, 1200, 0},
	roll{Dice{1: 0, 2: 2, 3: 2, 4: 2, 5: 0, 6: 0}, 0, 6},
}

func TestScore(t *testing.T) {
	for i, rl := range rolls {
		t.Logf("Run %v %+v:", i, rl.Dice)
		score, rem := Score(0, rl.Dice)
		if score != rl.Score {
			t.Errorf("    expected score %v, got %v", rl.Score, score)
		}
		if rem.N() != rl.Nrem {
			t.Errorf("    expected len(rem) %v, got %v", rl.Nrem, rem.N())
		}
	}
}

const ngames = 10000

func TestBreaker(t *testing.T) {
	var foo = [][]int{
		[]int{towin, 0, 0, 0},
		[]int{0, towin, 0, 1},
		[]int{0, 0, towin, 2},
		[]int{towin, towin, 0, 0},
		[]int{0, towin, towin, 1},
		[]int{towin, 0, towin, 0},
		[]int{towin, towin, towin, 0},
	}

	for _, tst := range foo {
		index := Breaker(tst[:3])
		if index != tst[3] {
			t.Errorf("test %+v failed: got index %v, expected %v", tst[:3], index, tst[3])
		}
	}
}

func TestPlayers(t *testing.T) {
	players := []Strategy{
		//GoForItStrategy(450),
		HoldStrategy{},
		HoldStrategy{},
		AggressiveEndStrategy{HoldStrategy{}},
	}

	rng := rand.New(rand.NewSource(time.Now().Unix()))

	counts := make([]int, len(players))
	for i := 0; i < ngames; i++ {
		scores := Play(rng, nil, players...)
		winner := Winner(scores)
		counts[winner]++
	}

	for i, v := range counts {
		t.Logf("counts[i]=%v", counts[i])
		t.Logf("Player %v won %.3f%% of matches", i, float64(v)/ngames*100)
	}
}
