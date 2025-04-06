package zen

import (
	"blackjack/pkg/strategy"
	"fmt"
	"math"
)

type Count struct {
	Count           int
	Cards           int
	InitialBestSize int
}

func (z *Count) RawCount() int {
	return z.Count
}

func NewZenCount(cards int) *Count {
	return &Count{
		Count:           0,
		Cards:           cards,
		InitialBestSize: 1,
	}
}

var cardValues = map[string]int{
	"Two":   1,
	"Three": 1,
	"Seven": 1,
	"Four":  2,
	"Five":  2,
	"Six":   2,
	"Eight": 0,
	"Nine":  0,
	"Ten":   -2,
	"Jack":  -2,
	"Queen": -2,
	"King":  -2,
	"Ace":   -2,
}

func (z *Count) ProcessCard(card string) {
	val, ok := cardValues[card]
	if ok {
		z.Count += val
		z.Cards -= 1
	}
}

func (z *Count) TrueCount() float64 {
	decksRemaining := float64(z.Cards) / 52
	roundedDecks := math.Ceil(decksRemaining)
	count := math.Abs(float64(z.Count)) / roundedDecks
	if z.Count < 0 {
		return count * -1
	}
	return count
}

func (z *Count) BetSize() int {
	count := z.TrueCount()
	countInt := int(count)
	if countInt > 6 {
		return z.InitialBestSize * 12
	} else if countInt > 5 {
		return z.InitialBestSize * 10
	} else if countInt > 4 {
		return z.InitialBestSize * 8
	} else if countInt > 3 {
		return z.InitialBestSize * 6
	} else if countInt > 2 {
		return z.InitialBestSize * 4
	} else if countInt > 1 {
		return z.InitialBestSize * 2
	}
	return z.InitialBestSize
}

func (z *Count) String() string {
	return fmt.Sprintf("True Count: %f, Cards remaining: %d, Bet Size: %d", z.TrueCount(), z.Cards, z.BetSize())
}

func (z *Count) EvaluateHand(playersCards []string, dealerCard string) strategy.Action {
	return strategy.Hit
}

func (z *Count) ShowCounts() map[string]int {
	return cardValues
}
