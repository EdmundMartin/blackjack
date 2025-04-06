package deck

type Hand struct {
	Cards    []*Card
	IsDealer bool
}

func (h *Hand) IsSoft() bool {
	return false
}
