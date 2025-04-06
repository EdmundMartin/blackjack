package deck

type Suit string

const (
	Hearts   Suit = "Hearts"
	Diamonds Suit = "Diamonds"
	Clubs    Suit = "Clubs"
	Spades   Suit = "Spades"
)

func GetSuits() []Suit {
	return []Suit{
		Hearts, Diamonds, Clubs, Spades,
	}
}

type Value string

const (
	Ace   Value = "Ace"
	Two   Value = "Two"
	Three Value = "Three"
	Four  Value = "Four"
	Five  Value = "Five"
	Six   Value = "Six"
	Seven Value = "Seven"
	Eight Value = "Eight"
	Nine  Value = "Nine"
	Ten   Value = "Ten"
	Jack  Value = "Jack"
	Queen Value = "Queen"
	King  Value = "King"
)

func GetCardValues() []Value {
	return []Value{
		Ace, Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King,
	}
}

func (v Value) NumericValue() int {
	switch v {
	case Ace:
		return 1
	case Two:
		return 2
	case Three:
		return 3
	case Four:
		return 4
	case Five:
		return 5
	case Six:
		return 6
	case Seven:
		return 7
	case Eight:
		return 8
	case Nine:
		return 9
	case Ten, Jack, Queen, King:
		return 10
	}
	return -1
}

type Card struct {
	Value Value
	Suit  Suit
}
