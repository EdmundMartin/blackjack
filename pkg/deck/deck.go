package deck

import (
	"blackjack/pkg/strategy"
	"blackjack/pkg/strategy/zen"
	"math/rand"
)

type Deck struct {
	Cards     []*Card
	DeckCount int
	Strategy  strategy.Strategy
}

func NewDeck(deckCount int) *Deck {
	cards := make([]*Card, 0, 52*deckCount)
	for i := 0; i < deckCount; i++ {
		for _, suit := range GetSuits() {
			for _, card := range GetCardValues() {
				cards = append(cards, &Card{
					Value: card,
					Suit:  suit,
				})
			}
		}
	}
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})

	zenCount := zen.NewZenCount(len(cards))
	return &Deck{
		Cards:     cards,
		DeckCount: deckCount,
		Strategy:  zenCount,
	}
}

func (d *Deck) DealCard() *Card {
	var x *Card
	x, d.Cards = d.Cards[len(d.Cards)-1], d.Cards[:len(d.Cards)-1]

	d.Strategy.ProcessCard(string(x.Value))
	return x
}

func (d *Deck) DealHand(isDealer bool) *Hand {
	leftCard := d.DealCard()
	rightCard := d.DealCard()
	return &Hand{
		Cards:    []*Card{leftCard, rightCard},
		IsDealer: isDealer,
	}
}
