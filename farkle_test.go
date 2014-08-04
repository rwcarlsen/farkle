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

const ngames = 100000

func TestPlayers(t *testing.T) {
	players := []Strategy{
		//GoForItStrategy(450),
		HoldStrategy{},
		HoldStrategy{},
	}

	rng := rand.New(rand.NewSource(time.Now().Unix()))

	goforit := 0.0
	hold := 0.0
	for i := 0; i < ngames; i++ {
		scores := Play(rng, nil, players...)
		winner := Winner(scores)
		if winner == 0 {
			goforit++
		} else {
			hold++
		}
	}
	t.Logf("GoForIt wins with %.3f%% of matches", goforit/ngames*100)
	t.Logf("Hold wins %.3f%% of matches", hold/ngames*100)
}
