package deck

import (
	"math/rand"
	"time"
)

type Card struct {
	Suit  string
	Value int
}

type Deck []Card

func AddCard(value int, suit string) Card {
	return Card{
		Value: value,
		Suit:  suit,
	}
}

func Shuffle(deck Deck) Deck {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	r.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	return deck
}

func NewDeck() Deck {
	suits := []string{"Spades", "Clubs", "Hearts", "Diamonds"}
	var deck Deck
	for _, suit := range suits {
		for i := 2; i <= 14; i++ {
			deck = append(deck, AddCard(i, suit))
		}
	}
	return deck
}
