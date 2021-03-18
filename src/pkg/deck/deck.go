package deck

import (
	"math/rand"
	"time"
)

type Card struct {
	Suit  string
	Value int
}

func (c Card) ToString() string {
	valuesMapping := map[int]string{
		2:  "2",
		3:  "3",
		4:  "4",
		5:  "5",
		6:  "6",
		7:  "7",
		8:  "8",
		9:  "9",
		10: "10",
		11: "J",
		12: "Q",
		13: "K",
		14: "A",
	}

	suitMapping := map[string]string{
		"Spades":   "♠",
		"Hearts":   "♥",
		"Diamonds": "♦",
		"Clubs":    "♣",
	}
	return suitMapping[c.Suit] + valuesMapping[c.Value]
}

func AddCard(value int, suit string) Card {
	return Card{
		Value: value,
		Suit:  suit,
	}
}

func NewDeck() []Card {
	suits := []string{"Spades", "Clubs", "Hearts", "Diamonds"}
	var deck []Card
	for _, suit := range suits {
		for i := 2; i <= 14; i++ {
			deck = append(deck, AddCard(i, suit))
		}
	}
	return deck
}

func Shuffle(deck []Card) []Card {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	r.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	return deck
}
