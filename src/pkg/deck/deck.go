package deck

import (
	"math/rand"
	"time"
)

type Card struct {
	Suit  string
	Value int
}

type Deck struct {
	Cards []Card
}

func AddCard(value int, suit string) Card {
	return Card{
		Value: value,
		Suit:  suit,
	}
}

func Shuffle(deck *Deck) *Deck {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	r.Shuffle(len(deck.Cards), func(i, j int) {
		deck.Cards[i], deck.Cards[j] = deck.Cards[j], deck.Cards[i]
	})
	return deck
}

func NewDeck() Deck {
	suits := []string{"Spades", "Clubs", "Hearts", "Diamonds"}
	var deck Deck
	for _, suit := range suits {
		for i := 2; i <= 14; i++ {
			deck.Cards = append(deck.Cards, AddCard(i, suit))
		}
	}
	return deck
}
