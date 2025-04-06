package basic

import (
	"blackjack/pkg/strategy"
	"fmt"
)

type BasicStrategy struct {
	Count int
	Size  int
	Cards int
}

func (b *BasicStrategy) ProcessCard(card string) {
	return
}

func (b *BasicStrategy) TrueCount() float64 {
	return float64(b.Count)
}

func (b *BasicStrategy) RawCount() int {
	return b.Count
}

func (b *BasicStrategy) BetSize() int {
	return b.Size
}

func (b *BasicStrategy) String() string {
	return fmt.Sprintf("True Count: %f, Cards remaining: %d, Bet Size: %d", b.TrueCount(), b.Cards, b.BetSize())
}

func (b *BasicStrategy) EvaluateHand(playersCards []string, dealerCard string) strategy.Action {
	return strategy.Hit
}
