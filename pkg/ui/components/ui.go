package components

import (
	"blackjack/pkg/game"
	"blackjack/pkg/ui"
	"fmt"
	"image/color"
	"path/filepath"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type UI struct {
	Game *game.Game

	Window      fyne.Window
	Balance     *widget.Label
	Status      *widget.Label
	BetEntry    *widget.Entry
	CountLabel  *widget.Label
	AdviceLabel *widget.Label

	ToggleCountBtn  *widget.Button
	ToggleAdviceBtn *widget.Button

	DealerBox *fyne.Container
	PlayerBox *fyne.Container
}

func NewUI(a fyne.App, g *game.Game) *UI {
	w := a.NewWindow("Blackjack")
	w.Resize(fyne.NewSize(700, 650))

	u := &UI{
		Game:        g,
		Window:      w,
		Balance:     widget.NewLabel("Balance: $1000"),
		Status:      widget.NewLabel("Place your bet to begin."),
		BetEntry:    widget.NewEntry(),
		CountLabel:  widget.NewLabel(""),
		AdviceLabel: widget.NewLabel(""),
		DealerBox:   container.NewHBox(),
		PlayerBox:   container.NewVBox(),
	}

	u.BetEntry.SetPlaceHolder("Bet Amount")
	u.BetEntry.SetText("1")

	u.ToggleCountBtn = widget.NewButton("Show Count", func() {
		u.Game.ShowCount = !u.Game.ShowCount
		if u.Game.ShowCount {
			u.ToggleCountBtn.SetText("Hide Count")
		} else {
			u.ToggleCountBtn.SetText("Show Count")
			u.CountLabel.SetText("")
		}
		u.Game.UpdateCounts()
	})

	u.ToggleAdviceBtn = widget.NewButton("Show Count Advice", func() {
		u.Game.ShowAdvice = !u.Game.ShowAdvice
		if u.Game.ShowAdvice {
			u.ToggleAdviceBtn.SetText("Hide Count Advice")
		} else {
			u.ToggleAdviceBtn.SetText("Show Count Advice")
			u.AdviceLabel.SetText("")
		}
		u.Game.UpdateCounts()
	})

	placeBetButton := widget.NewButton("Place Bet", func() {
		bet, _ := strconv.Atoi(u.BetEntry.Text)
		u.Game.PlaceBet(bet)
	})

	buttonRow := container.NewHBox(
		widget.NewButton("Hit", func() { u.Game.Hit() }),
		widget.NewButton("Stand", func() { u.Game.Stand() }),
		widget.NewButton("Double", func() { u.Game.DoubleDown() }),
		widget.NewButton("Split", func() { u.Game.SplitHand() }),
	)

	green := canvas.NewRectangle(&color.RGBA{R: 0, G: 100, B: 0, A: 255})
	content := container.NewVBox(
		u.Balance,
		container.NewHBox(u.BetEntry, placeBetButton),
		u.ToggleCountBtn,
		u.CountLabel,
		u.ToggleAdviceBtn,
		u.AdviceLabel,
		widget.NewLabel("Dealer:"),
		u.DealerBox,
		widget.NewLabel("Player:"),
		u.PlayerBox,
		buttonRow,
		u.Status,
	)

	mainContainer := container.NewMax(green, container.NewPadded(content))
	u.Window.SetContent(mainContainer)
	return u
}

func (u *UI) Render() {
	u.PlayerBox.Objects = nil
	for i, hand := range u.Game.PlayerHands {
		row := container.NewHBox()
		for _, card := range hand.Hand.Cards {
			img := canvas.NewImageFromFile(ui.CardImagePath(card))
			img.SetMinSize(fyne.NewSize(60, 90))
			row.Add(img)
		}
		label := widget.NewLabel(fmt.Sprintf("Hand %d", i+1))
		u.PlayerBox.Add(container.NewVBox(label, row))
	}

	row := container.NewHBox()
	for i, card := range u.Game.DealerHand.Cards {
		img := canvas.NewImageFromFile(func() string {
			if i == 0 && u.Game.DealerHidden {
				return filepath.Join("assets", "back.png")
			}
			return ui.CardImagePath(card)
		}())
		img.SetMinSize(fyne.NewSize(60, 90))
		row.Add(img)
	}
	u.DealerBox.Objects = row.Objects
	u.PlayerBox.Refresh()
	u.DealerBox.Refresh()
}
