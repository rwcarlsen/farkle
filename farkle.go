package farkle

import (
	"math/rand"
	"time"
)

const (
	towin  = 5000
	thresh = 350
)

type Context struct {
	// AllScores contains the scores of all players (including this
	// player).
	AllScores []int
	// YourScore is the player's current total score excluding unscored points
	// accumulated in the current turn.
	YourScore int
	// Points represents the points accumulated from prior rolls on the
	// currently active turn on - points that are not permanent and may be
	// lose.
	Points int
	// EndScore is the number of points a player must have to end the game.
	EndScore int
	// TurnThresh is the minimum number of points that must be accumulated on a
	// turn for them to become a permanent part of a player's score.
	TurnThresh int
}

type Strategy interface {
	// Roll represents a roll of the dice.  keep is the dice values that are
	// set aside and not rolled.  Returning a non-nil keep slice results in
	// another roll for this turn.  Returning a nil keep slice ends the turn.
	Roll(c Context, got []int) (keep []int)
}

type Game struct {
	Rng     *rand.Rand
	players []Strategy
	scores  []int
}

func (g *Game) AddPlayer(s Strategy) {
	g.players = append(g.players, s)
	g.scores = append(g.scores, 0)
}

func (g *Game) lead() (index int) {
	best := 0
	for i, v := range g.scores {
		if v > best {
			best = v
			index = i
		}
	}
	return index
}

func (g *Game) turn(s Strategy, score int) (points int) {
	c := Context{
		AllScores:  append([]int{}, scores...),
		YourScore:  score,
		Points:     0,
		EndScore:   towin,
		TurnThresh: thresh,
	}

	ndice := 5
	for ndice > 0 {
		got := rolldice(ndice)
		keep := s.Roll(c, got)
	}
}

type Dice map[int]int

func (d Dice) Score() int {
}

func (d Dice) Clone() Dice {
	clone := Dice{}
	for x, n := range d {
		clone[x] = n
	}
	return clone
}

func (d Dice) Sub(other Dice) Dice {
	diff := Dice{}
	for x, n := range d {
		if nother := other[x]; nother > 0 && n > nother {
			diff[x] = n - nother
		}
	}
	return diff
}

func scoreStraight(d Dice) (score int, rem Dice) {
	for i := 1; i <= 6; i++ {
		if d[i] == 0 {
			return 0, d.Clone()
		}
	}
	return 1000, Dice{}
}

func scoreOneFive(d Dice) (score int, rem Dice) {
	rem = d.Clone()
	if 
	for i := 1; i <= 6; i++ {
		if d[i] == 1 {
			return 0, Dice{}
		}
	}
	return 0, d.Clone()
}

func scoreTriple(d Dice) (score int, rem Dice) {
	rem = d.Clone()
	for x, n := range d {
		if n >= 3 {
			if x == 1 {
				score += 1000
			} else {
				score += x * 100
			}
			rem[x] -= 3
		}
	}
	return score, rem
}

func (g *Game) Run() {
	if g.Rng == nil {
		g.Rng = rand.New(rand.NewSource(time.Now().Unix()))
	}

	for g.lead() < towin {
		for i, p := range g.players {
			c := scores
		}
	}

}

func rolldice(rng *rand.Rand, n int) []int {
	dice := make([]int, n)
	for i := range dice {
		dice[i] = rng.Intn(6)
	}
	return dice
}
