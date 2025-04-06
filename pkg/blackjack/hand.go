package blackjack

import "blackjack/pkg/deck"

func HandValue(h *deck.Hand) (int, bool) {
	total := 0
	aces := 0

	for _, c := range h.Cards {
		v := c.Value.NumericValue()
		total += v
		if c.Value == deck.Ace {
			aces++
		}
	}

	isSoft := false
	for aces > 0 && total+10 <= 21 {
		total += 10
		aces--
		isSoft = true
	}

	return total, isSoft
}

func IsBlackjack(h *deck.Hand) bool {
	val, _ := HandValue(h)
	return val == 21 && len(h.Cards) == 2
}

func IsBust(h *deck.Hand) bool {
	val, _ := HandValue(h)
	return val > 21
}

func CanSplit(h *deck.Hand) bool {
	if len(h.Cards) != 2 {
		return false
	}
	return h.Cards[0].Value == h.Cards[1].Value
}
