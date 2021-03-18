package deck

import (
	"fmt"
	"reflect"
	"testing"
)

func TestDeck(t *testing.T) {
	t.Run("must return 52 cards", func(t *testing.T) {
		deck := NewDeck()
		if len(deck) < 52 {
			t.Errorf("want %d got %d cards", 52, len(deck))
		}
	})
}

func TestShuffle(t *testing.T) {
	t.Run("must return shuffled deck", func(t *testing.T) {
		deck := NewDeck()
		deck2 := NewDeck()
		first := Shuffle(deck)
		second := Shuffle(deck2)
		if reflect.DeepEqual(first, second) {
			t.Errorf("the first card still on first position")
		}
	})
}

func TestAddCard(t *testing.T) {
	t.Run("must return 1 card", func(t *testing.T) {
		card := AddCard(2, "Spades")
		fmt.Printf("%v", card)
		if card.Value != 2 || card.Suit != "Spades" {
			t.Errorf("want 2 Of Spades, got %v", card)
		}
	})
}

func TestCardToString(t *testing.T) {
	t.Run("must return ♠2", func(t *testing.T) {
		card := Card{Value: 2, Suit: "Spades"}.ToString()
		if card != "♠2" {
			t.Errorf("want ♠2 got %v", card)
		}
	})
}
