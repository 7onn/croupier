package deck

import (
	"fmt"
	"reflect"
	"testing"
)

func TestDeck(t *testing.T) {
	t.Run("must return 52 cards", func(t *testing.T) {
		deck := NewDeck()
		if len(deck.Cards) < 52 {
			t.Errorf("want %d got %d cards", 52, len(deck.Cards))
		}
	})
}

func TestShuffle(t *testing.T) {
	t.Run("must return shuffled deck", func(t *testing.T) {
		deck := NewDeck()
		deck2 := NewDeck()
		first := Shuffle(&deck)
		second := Shuffle(&deck2)
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
