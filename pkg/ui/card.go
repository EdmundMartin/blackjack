package ui

import (
	"blackjack/pkg/deck"
	"fmt"
	"path/filepath"
	"strings"
)

func CardImagePath(card *deck.Card) string {
	if card == nil {
		return "assets/back.png"
	}

	val := strings.ToUpper(string(card.Value)[0:1])
	switch card.Value {
	case deck.Ten:
		val = "10"
	case deck.Jack:
		val = "J"
	case deck.Queen:
		val = "Q"
	case deck.King:
		val = "K"
	case deck.Ace:
		val = "A"
	case deck.Two, deck.Three, deck.Four, deck.Five, deck.Six,
		deck.Seven, deck.Eight, deck.Nine:
		val = fmt.Sprintf("%d", card.Value.NumericValue())
	}

	var suit string
	switch card.Suit {
	case deck.Hearts:
		suit = "H"
	case deck.Diamonds:
		suit = "D"
	case deck.Clubs:
		suit = "C"
	case deck.Spades:
		suit = "S"
	}

	filename := fmt.Sprintf("%s%s.png", val, suit)
	return filepath.Join("assets", filename)
}
