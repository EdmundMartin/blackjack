package main

import (
	"blackjack/pkg/blackjack"
	"blackjack/pkg/deck"
	"blackjack/pkg/ui"
	"fmt"
	"image/color"
	"math/rand"
	"path/filepath"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	gameDeck         *deck.Deck
	dealerHand       *deck.Hand
	playerHands      []*blackjack.GameHand
	currentHandIndex int
	gameBalance      = 1000
	dealerHidden     = true

	// UI elements
	balanceLabel    *widget.Label
	statusLabel     *widget.Label
	betEntry        *widget.Entry
	playerBox       *fyne.Container
	dealerBox       *fyne.Container
	hitBtn          *widget.Button
	standBtn        *widget.Button
	doubleBtn       *widget.Button
	splitBtn        *widget.Button
	placeBetButton  *widget.Button
	countLabel      *widget.Label
	toggleCountBtn  *widget.Button
	showCount                                = false
	loadedImages    map[string]*canvas.Image = map[string]*canvas.Image{}
	showAdvice                               = false
	adviceLabel     *widget.Label
	toggleAdviceBtn *widget.Button
)

func toggleAdvice() {
	showAdvice = !showAdvice
	if showAdvice {
		toggleAdviceBtn.SetText("Hide Count Advice")
		updateAdviceLabel()
	} else {
		toggleAdviceBtn.SetText("Show Count Advice")
		adviceLabel.SetText("")
	}
}

func updateAdviceLabel() {
	if !showAdvice || gameDeck == nil {
		adviceLabel.SetText("")
		return
	}

	counts := gameDeck.Strategy.ShowCounts()
	order := []string{
		"Ace", "Two", "Three", "Four", "Five", "Six",
		"Seven", "Eight", "Nine", "Ten", "Jack", "Queen", "King",
	}

	output := "Count Advice:\n"
	for _, val := range order {
		c, ok := counts[val]
		if !ok {
			c = 0
		}
		output += fmt.Sprintf("%-6s: %+d\n", val, c)
	}
	adviceLabel.SetText(output)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	a := app.New()
	w := a.NewWindow("Blackjack - Full Game")
	w.Resize(fyne.NewSize(700, 650))

	balanceLabel = widget.NewLabel("Balance: $1000")
	statusLabel = widget.NewLabel("Place your bet to begin.")
	betEntry = widget.NewEntry()
	betEntry.SetPlaceHolder("Bet Amount")
	betEntry.SetText("1")

	playerBox = container.NewVBox()
	dealerBox = container.NewHBox()

	countLabel = widget.NewLabel("")
	toggleCountBtn = widget.NewButton("Show Count", toggleCount)

	placeBetButton = widget.NewButton("Place Bet", placeBet)
	hitBtn = widget.NewButton("Hit", hit)
	standBtn = widget.NewButton("Stand", stand)
	doubleBtn = widget.NewButton("Double Down", doubleDown)
	splitBtn = widget.NewButton("Split", splitHand)

	adviceLabel = widget.NewLabel("")
	toggleAdviceBtn = widget.NewButton("Show Count Advice", toggleAdvice)

	buttons := container.NewHBox(hitBtn, standBtn, doubleBtn, splitBtn)
	betUI := container.NewHBox(betEntry, placeBetButton)

	casinoGreen := &color.RGBA{R: 0, G: 100, B: 0, A: 255} // rich green
	bg := canvas.NewRectangle(casinoGreen)

	content := container.NewVBox(
		balanceLabel,
		betUI,
		toggleCountBtn,
		countLabel,
		toggleAdviceBtn,
		adviceLabel,
		widget.NewLabel("Dealer:"),
		dealerBox,
		widget.NewLabel("Player:"),
		playerBox,
		buttons,
		statusLabel,
	)

	mainContent := container.NewMax(bg, container.NewPadded(content))

	toggleActionButtons(false)
	w.SetContent(mainContent)
	w.ShowAndRun()
}

func placeBet() {
	if len(playerHands) > 0 {
		statusLabel.SetText("Finish current round first.")
		return
	}
	bet, err := strconv.Atoi(betEntry.Text)
	if err != nil || bet <= 0 {
		statusLabel.SetText("Enter a valid bet amount.")
		return
	}
	if bet > gameBalance {
		statusLabel.SetText("Not enough balance.")
		return
	}

	gameDeck = deck.NewDeck(6)
	dealerHand = gameDeck.DealHand(true)
	playerHand := gameDeck.DealHand(false)

	dealerHidden = true

	if dealerHand.Cards[1].Value == deck.Ace {
		statusLabel.SetText("Dealer shows Ace. Insurance not implemented yet.")
	}

	gameBalance -= bet
	balanceLabel.SetText("Balance: $" + strconv.Itoa(gameBalance))
	playerHands = []*blackjack.GameHand{
		{Hand: playerHand, Bet: bet, Active: true},
	}
	currentHandIndex = 0
	updateHands()
	updateCountLabel()
	statusLabel.SetText("Your move.")
	toggleActionButtons(true)
}

func hit() {
	if currentHandIndex >= len(playerHands) {
		return
	}
	hand := playerHands[currentHandIndex]
	hand.Hand.Cards = append(hand.Hand.Cards, gameDeck.DealCard())
	updateHands()
	updateCountLabel()

	if blackjack.IsBust(hand.Hand) {
		statusLabel.SetText("Busted!")
		nextHand()
	}
}

func stand() {
	playerHands[currentHandIndex].Active = false
	statusLabel.SetText("Standing.")
	nextHand()
}

func doubleDown() {
	hand := playerHands[currentHandIndex]
	if gameBalance < hand.Bet {
		statusLabel.SetText("Not enough to double down.")
		return
	}
	gameBalance -= hand.Bet
	balanceLabel.SetText("Balance: $" + strconv.Itoa(gameBalance))
	blackjack.DoubleDown(hand, gameDeck)
	updateHands()
	updateCountLabel()
	statusLabel.SetText("Doubled down.")
	nextHand()
}

func splitHand() {
	hand := playerHands[currentHandIndex]
	if !blackjack.CanSplit(hand.Hand) {
		statusLabel.SetText("Can't split.")
		return
	}
	if gameBalance < hand.Bet {
		statusLabel.SetText("Not enough to split.")
		return
	}
	gameBalance -= hand.Bet
	balanceLabel.SetText("Balance: $" + strconv.Itoa(gameBalance))

	newHands := blackjack.SplitHand(hand, gameDeck, hand.Bet)
	playerHands = append(playerHands[:currentHandIndex], append(newHands, playerHands[currentHandIndex+1:]...)...)
	updateHands()
	updateCountLabel()
	statusLabel.SetText("Hand split.")
}

func nextHand() {
	currentHandIndex++
	if currentHandIndex >= len(playerHands) {
		dealerHidden = false
		blackjack.DealerPlay(dealerHand, gameDeck)
		updateHands()
		updateCountLabel()

		for _, hand := range playerHands {
			result := blackjack.EvaluateHand(hand, dealerHand)
			switch result {
			case blackjack.BlackjackWin:
				gameBalance += int(float64(hand.Bet) * 2.5)
			case blackjack.Win:
				gameBalance += hand.Bet * 2
			case blackjack.Push:
				gameBalance += hand.Bet
			case blackjack.Loss:
			}
		}
		statusLabel.SetText("Round complete. Place bet to play again.")
		playerHands = nil
		currentHandIndex = 0
		balanceLabel.SetText("Balance: $" + strconv.Itoa(gameBalance))
		toggleActionButtons(false)
		return
	}
	updateHands()
	updateCountLabel()
	statusLabel.SetText("Next hand.")
}

func updateHands() {
	playerBox.Objects = nil

	for i, hand := range playerHands {
		cards := handToImages(hand.Hand, false)
		label := widget.NewLabel("Hand " + strconv.Itoa(i+1))
		row := container.NewVBox(label, container.NewHBox(cards...))
		playerBox.Add(row)
	}

	dealerBox.Objects = handToImages(dealerHand, dealerHidden)
	playerBox.Refresh()
	dealerBox.Refresh()
}

func handToImages(h *deck.Hand, hideFirst bool) []fyne.CanvasObject {
	var imgs []fyne.CanvasObject
	for i, c := range h.Cards {
		var img *canvas.Image
		var path string
		if hideFirst && i == 0 {
			path = filepath.Join("assets", "back.png")
		} else {
			path = ui.CardImagePath(c)
		}
		img, ok := loadedImages[path]

		if !ok {
			fmt.Println(path)
			img = canvas.NewImageFromFile(path)
			loadedImages[path] = img
		}
		img.SetMinSize(fyne.NewSize(60, 90))
		imgs = append(imgs, img)
	}
	return imgs
}

func toggleActionButtons(enable bool) {
	hitBtn.Disable()
	standBtn.Disable()
	doubleBtn.Disable()
	splitBtn.Disable()

	if enable {
		hitBtn.Enable()
		standBtn.Enable()
		doubleBtn.Enable()
		splitBtn.Enable()
	}
}

func toggleCount() {
	showCount = !showCount
	if showCount {
		toggleCountBtn.SetText("Hide Count")
		updateCountLabel()
	} else {
		toggleCountBtn.SetText("Show Count")
		countLabel.SetText("")
	}
}

func updateCountLabel() {
	if gameDeck == nil || !showCount {
		countLabel.SetText("")
		return
	}
	raw := gameDeck.Strategy.RawCount()
	trueCount := gameDeck.Strategy.TrueCount()

	countLabel.SetText(fmt.Sprintf("Raw Count: %+d\nTrue Count: %.2f", raw, trueCount))
}
