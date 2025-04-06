package strategy

type Action string

const (
	Hit       Action = "Hit"
	Stand     Action = "Stand"
	Split     Action = "Split"
	Double    Action = "Double"
	Insurance Action = "Insurance"
)

type Strategy interface {
	ProcessCard(card string)
	TrueCount() float64
	RawCount() int
	BetSize() int
	String() string
	EvaluateHand(playersCards []string, dealerCard string) Action
	ShowCounts() map[string]int
}
