package farkle

import (
	"math/rand"
	"time"
)

const (
	towin  = 5000
	thresh = 350
	ndice  = 6
)

type Context struct {
	// scores contains the scores of all players (including this
	// player).
	scores []int
	// Index is the player's index - useful for retrieving the player's score.
	Index int
	// Points represents the points accumulated from prior rolls on the
	// currently active turn on - points that are not permanent and may be
	// lose.
	Points int
	// EndScore is the number of points a player must have to end the game.
	EndScore int
	// TurnThresh is the minimum number of points that must be accumulated on a
	// turn for them to become a permanent part of a player's score.
	TurnThresh int
	// ScoreFn is the function used to calculate all scores from the dice.
	ScoreFn ScoreFunc
	// Last is true if this is the game's last turn/round
	Last bool
}

func (c Context) Score(i int) int { return c.scores[i] }

type Strategy interface {
	// Roll represents a roll of the dice.  keep is the dice values that are
	// set aside and not rolled.  Returning nil scores all possible points and
	// ends the turn.
	Roll(c Context, got Dice) (keep Dice)
}

type HoldStrategy struct{}

func (_ HoldStrategy) Roll(c Context, got Dice) (keep Dice) {
	var pts int
	pts, keep = KeepMax(c.ScoreFn, got)
	if c.Points+pts >= c.TurnThresh {
		return nil
	}
	return keep
}

func validKeep(got, keep Dice) bool {
	for i, n := range keep {
		if n > got[i] {
			return false
		}
	}
	return true
}

func Turn(ctx Context, rng *rand.Rand, s Strategy) (points int) {
	var d, rem Dice
	var pts int
	n := ndice
	for n > 0 {
		ctx.Points = points
		d = RollDice(rng, n, d)

		// check for failure to roll scoring dice
		if pts, _ := ctx.ScoreFn(0, d); pts == 0 {
			return 0
		}

		keep := s.Roll(ctx, d)

		// check for nil return indicating turn termination
		if keep == nil {
			points, _ = ctx.ScoreFn(points, d)
			return points
		}

		pts, rem = ctx.ScoreFn(0, keep)
		points += pts
		n -= keep.N()

		if pts == 0 {
			panic("keep dice don't score")
		} else if !validKeep(d, keep) {
			panic("keep dice are not a subset of got dice")
		}

		// check for hot dice
		if n == 0 && rem.N() == 0 {
			n = ndice
		}
	}
	return points
}

type Dice []int

func NewDice() Dice {
	return make(Dice, 7) // num die sides plus one
}

func (d Dice) N() int {
	tot := 0
	for _, n := range d {
		tot += n
	}
	return tot
}

func (d Dice) Clone() Dice {
	clone := make(Dice, len(d))
	for i, n := range d {
		clone[i] = n
	}
	return clone
}

func Play(rng *rand.Rand, fn ScoreFunc, players ...Strategy) (scores []int) {
	scores = make([]int, len(players))
	if rng == nil {
		rng = rand.New(rand.NewSource(time.Now().Unix()))
	}
	if fn == nil {
		fn = Score
	}

	done := false
	for !done {
		// give one turn to all players after one breaks winner threshold
		done = Breaker(scores) >= 0

		for i, p := range players {
			ctx := Context{
				scores:     scores,
				ScoreFn:    fn,
				Index:      i,
				EndScore:   towin,
				TurnThresh: thresh,
				Last:       done,
			}
			scores[i] += Turn(ctx, rng, p)
		}

		// keep going if there is a tie
		if len(Winners(scores)) > 1 {
			done = false
		}
	}

	return scores
}

// Breaker returns the index of the first player who broke the game-end
// threshold or -1 if no player has broken it yet.
func Breaker(scores []int) (index int) {
	for i, v := range scores {
		if v >= towin {
			return i
		}
	}
	return -1
}

// Winner returns the index of the last player with the highest score.
func Winner(scores []int) (index int) {
	best := 0
	for i, v := range scores {
		if v > best {
			index = i
			best = v
		}
	}
	return index
}

// Winners returns the indices of the players with the highest score.
func Winners(scores []int) (indices []int) {
	best := 0
	for _, v := range scores {
		if v > best {
			best = v
		}
	}
	for i, v := range scores {
		if v == best {
			indices = append(indices, i)
		}
	}
	return indices
}

func RollDice(rng *rand.Rand, n int, d Dice) Dice {
	if d == nil {
		d = NewDice()
	}
	for i := 0; i < n; i++ {
		d[rng.Intn(6)+1]++
	}
	return d
}

func KeepMax(fn ScoreFunc, d Dice) (points int, scoring Dice) {
	points, rem := fn(0, d)
	scoring = d.Clone()
	for x, n := range rem {
		scoring[x] -= n
	}
	return points, scoring
}

// ScoreFunc returns the highest score possible for the given dice.
type ScoreFunc func(prevscore int, d Dice) (score int, rem Dice)

func Score(prevscore int, d Dice) (score int, rem Dice) {
	return scoreOneFive(
		scoreTriple(
			scoreStraight(prevscore, d.Clone()),
		),
	)
}

func scoreStraight(prevscore int, d Dice) (score int, rem Dice) {
	for _, n := range d {
		if n == 0 {
			return prevscore, d
		}
	}
	return prevscore + 1500, NewDice()
}

func scoreOneFive(prevscore int, d Dice) (score int, rem Dice) {
	d[1] = 0
	d[5] = 0
	return prevscore + d[1]*100 + d[5]*50, d
}

func scoreTriple(prevscore int, d Dice) (score int, rem Dice) {
	for i, n := range d {
		if n >= 3 {
			if i == 1 {
				score += 1000 * (n / 3)
			} else {
				score += i * 100 * (n / 3)
			}
			d[i] = n % 3
		}
	}
	return prevscore + score, d
}
