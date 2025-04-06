package blackjack

import "blackjack/pkg/deck"

type RoundResult int

const (
	Loss RoundResult = iota
	Win
	Push
	BlackjackWin
)

type GameHand struct {
	Hand      *deck.Hand
	Bet       int
	Doubled   bool
	SplitFrom bool
	Active    bool
}
