package game

import (
	"blackjack/pkg/blackjack"
	"blackjack/pkg/deck"
	"fmt"
)

type Game struct {
	Deck          *deck.Deck
	DealerHand    *deck.Hand
	PlayerHands   []*blackjack.GameHand
	CurrentHand   int
	Balance       int
	DealerHidden  bool
	ShowCount     bool
	ShowAdvice    bool
	UpdateUI      func()
	UpdateStatus  func(string)
	UpdateBalance func(int)
	UpdateCount   func(raw int, true float64)
	UpdateAdvice  func(advice string)
}

func (g *Game) checkState() bool {
	if len(g.PlayerHands) == 0 || g.CurrentHand >= len(g.PlayerHands) {
		g.UpdateStatus("No active hand. Place a bet to start.")
		return true
	}
	return false
}

func NewGame(updateUI func(), updateStatus func(string), updateBalance func(int), updateCount func(int, float64), updateAdvice func(string)) *Game {
	return &Game{
		Balance:       1000,
		DealerHidden:  true,
		UpdateUI:      updateUI,
		UpdateStatus:  updateStatus,
		UpdateBalance: updateBalance,
		UpdateCount:   updateCount,
		UpdateAdvice:  updateAdvice,
	}
}

func (g *Game) PlaceBet(bet int) {
	if len(g.PlayerHands) > 0 {
		g.UpdateStatus("Finish current round first.")
		return
	}
	if bet <= 0 || bet > g.Balance {
		g.UpdateStatus("Invalid or insufficient bet.")
		return
	}

	g.Deck = deck.NewDeck(6)
	g.DealerHand = g.Deck.DealHand(true)
	playerHand := g.Deck.DealHand(false)
	g.DealerHidden = true

	g.Balance -= bet
	g.UpdateBalance(g.Balance)

	g.PlayerHands = []*blackjack.GameHand{{Hand: playerHand, Bet: bet, Active: true}}
	g.CurrentHand = 0

	if blackjack.IsBlackjack(playerHand) {
		if blackjack.IsBlackjack(g.DealerHand) {
			g.Balance += bet // push
			g.UpdateBalance(g.Balance)
			g.UpdateStatus("Push — both have blackjack.")
		} else {
			win := int(float64(bet) * 2.5)
			g.Balance += win
			g.UpdateBalance(g.Balance)
			g.UpdateStatus(fmt.Sprintf("Blackjack! You win $%d", win))
		}
		g.PlayerHands = nil
		g.CurrentHand = 0
		g.DealerHidden = false
		g.UpdateUI()
		g.UpdateCounts()
		return
	}

	g.UpdateUI()
	g.UpdateCounts()
	g.UpdateStatus("Your move.")
}

func (g *Game) Hit() {
	terminate := g.checkState()
	if terminate {
		return
	}

	hand := g.PlayerHands[g.CurrentHand]
	hand.Hand.Cards = append(hand.Hand.Cards, g.Deck.DealCard())
	g.UpdateUI()
	g.UpdateCounts()
	if blackjack.IsBust(hand.Hand) {
		g.UpdateStatus("Busted!")
		g.NextHand()
	}
}

func (g *Game) Stand() {
	terminate := g.checkState()
	if terminate {
		return
	}

	g.PlayerHands[g.CurrentHand].Active = false
	g.UpdateStatus("Standing.")
	g.NextHand()
}

func (g *Game) DoubleDown() {
	terminate := g.checkState()
	if terminate {
		return
	}

	hand := g.PlayerHands[g.CurrentHand]
	if g.Balance < hand.Bet {
		g.UpdateStatus("Not enough to double down.")
		return
	}
	g.Balance -= hand.Bet
	g.UpdateBalance(g.Balance)

	blackjack.DoubleDown(hand, g.Deck)
	g.UpdateUI()
	g.UpdateCounts()
	g.UpdateStatus("Doubled down.")
	g.NextHand()
}

func (g *Game) SplitHand() {
	terminate := g.checkState()
	if terminate {
		return
	}

	hand := g.PlayerHands[g.CurrentHand]
	if !blackjack.CanSplit(hand.Hand) {
		g.UpdateStatus("Can't split.")
		return
	}
	if g.Balance < hand.Bet {
		g.UpdateStatus("Not enough to split.")
		return
	}
	g.Balance -= hand.Bet
	g.UpdateBalance(g.Balance)

	newHands := blackjack.SplitHand(hand, g.Deck, hand.Bet)
	g.PlayerHands = append(g.PlayerHands[:g.CurrentHand], append(newHands, g.PlayerHands[g.CurrentHand+1:]...)...)
	g.UpdateUI()
	g.UpdateCounts()
	g.UpdateStatus("Hand split.")
}

func (g *Game) NextHand() {
	g.CurrentHand++
	if g.CurrentHand >= len(g.PlayerHands) {
		g.DealerHidden = false
		blackjack.DealerPlay(g.DealerHand, g.Deck)
		g.UpdateUI()
		g.UpdateCounts()

		for _, hand := range g.PlayerHands {
			result := blackjack.EvaluateHand(hand, g.DealerHand)
			switch result {
			case blackjack.BlackjackWin:
				g.Balance += int(float64(hand.Bet) * 2.5)
			case blackjack.Win:
				g.Balance += hand.Bet * 2
			case blackjack.Push:
				g.Balance += hand.Bet
			}
		}

		summary := "Round complete:\n"

		for i, hand := range g.PlayerHands {
			result := blackjack.EvaluateHand(hand, g.DealerHand)
			outcome := ""
			switch result {
			case blackjack.BlackjackWin:
				win := int(float64(hand.Bet) * 2.5)
				g.Balance += win
				outcome = fmt.Sprintf("Blackjack! ($%d)", win)
			case blackjack.Win:
				win := hand.Bet * 2
				g.Balance += win
				outcome = fmt.Sprintf("Win ($%d)", win)
			case blackjack.Push:
				g.Balance += hand.Bet
				outcome = "Push"
			case blackjack.Loss:
				outcome = "Loss"
			}
			summary += fmt.Sprintf("Hand %d: %s\n", i+1, outcome)
		}
		g.UpdateStatus(summary)
		g.PlayerHands = nil
		g.CurrentHand = 0
		g.UpdateBalance(g.Balance)
		return
	}
	g.UpdateUI()
	g.UpdateCounts()
	g.UpdateStatus("Next hand.")
}

func (g *Game) UpdateCounts() {
	if g.Deck == nil {
		return
	}

	if g.ShowCount {
		raw := g.Deck.Strategy.RawCount()
		trueCount := g.Deck.Strategy.TrueCount()
		g.UpdateCount(raw, trueCount)
	}
	if g.ShowAdvice {
		counts := g.Deck.Strategy.ShowCounts()
		order := []string{
			"Ace", "Two", "Three", "Four", "Five", "Six",
			"Seven", "Eight", "Nine", "Ten", "Jack", "Queen", "King",
		}
		advice := "Count Advice:\n"
		for _, val := range order {
			c := counts[val]
			advice += fmt.Sprintf("%-6s: %+d\n", val, c)
		}
		g.UpdateAdvice(advice)
	}
}
