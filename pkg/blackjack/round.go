package blackjack

import "blackjack/pkg/deck"

func DealerPlay(dealer *deck.Hand, d *deck.Deck) {
	for {
		val, soft := HandValue(dealer)
		if val > 21 {
			return
		}
		if val > 17 || (val == 17 && !soft) {
			return
		}
		dealer.Cards = append(dealer.Cards, d.DealCard())
	}
}

func EvaluateHand(player *GameHand, dealer *deck.Hand) RoundResult {
	if IsBlackjack(player.Hand) && !IsBlackjack(dealer) {
		return BlackjackWin
	}
	if IsBlackjack(dealer) && !IsBlackjack(player.Hand) {
		return Loss
	}

	playerVal, _ := HandValue(player.Hand)
	dealerVal, _ := HandValue(dealer)

	if playerVal > 21 {
		return Loss
	}
	if dealerVal > 21 {
		return Win
	}
	if playerVal > dealerVal {
		return Win
	}
	if playerVal < dealerVal {
		return Loss
	}
	return Push
}

func SplitHand(original *GameHand, d *deck.Deck, bet int) []*GameHand {
	c1 := original.Hand.Cards[0]
	c2 := original.Hand.Cards[1]

	h1 := &deck.Hand{Cards: []*deck.Card{c1, d.DealCard()}}
	h2 := &deck.Hand{Cards: []*deck.Card{c2, d.DealCard()}}

	return []*GameHand{
		{Hand: h1, Bet: bet, SplitFrom: true, Active: true},
		{Hand: h2, Bet: bet, SplitFrom: true, Active: true},
	}
}

func DoubleDown(player *GameHand, d *deck.Deck) {
	if player.Doubled || len(player.Hand.Cards) != 2 {
		return
	}
	player.Hand.Cards = append(player.Hand.Cards, d.DealCard())
	player.Doubled = true
	player.Bet *= 2
	player.Active = false
}
