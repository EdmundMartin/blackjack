package main

import (
	"blackjack/pkg/game"
	"blackjack/pkg/ui/components"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	var uiRef *components.UI

	blackJackGame := game.NewGame(
		func() { uiRef.Render() },
		func(msg string) { uiRef.Status.SetText(msg) },
		func(balance int) { uiRef.Balance.SetText("Balance: $" + strconv.Itoa(balance)) },
		func(raw int, trueCount float64) {
			uiRef.CountLabel.SetText(fmt.Sprintf("Raw: %+d \nTrue: %.2f", raw, trueCount))
		},
		func(advice string) {
			uiRef.AdviceLabel.SetText(advice)
		},
	)

	uiRef = components.NewUI(a, blackJackGame)
	uiRef.Window.ShowAndRun()
}
