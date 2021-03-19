package cards

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
		if card.Rank != 2 || card.Suit != "Spades" {
			t.Errorf("want 2 Of Spades, got %v", card)
		}
	})
}

func TestCardToString(t *testing.T) {
	t.Run("must return ♠2", func(t *testing.T) {
		card := Card{Rank: 2, Suit: "Spades"}.ToString()
		if card != "♠2" {
			t.Errorf("want ♠2 got %v", card)
		}
	})
}

func TestCopyDeck(t *testing.T) {
	t.Run("must return identical set of cards", func(t *testing.T) {
		deck := NewDeck()
		copied := CopyDeck(deck)
		if !reflect.DeepEqual(deck, copied) {
			t.Errorf("want %v got %v", deck, copied)
		}
	})
}

func TestSortCards(t *testing.T) {
	t.Run("must return sorted cards", func(t *testing.T) {
		cards := NewDeck()
		shuffled := Shuffle(cards)
		sorted := SortCards(shuffled)
		if !reflect.DeepEqual(cards, sorted) {
			t.Errorf("want %v got %v", cards, sorted)
		}
	})
}

func TestHandScoreToString(t *testing.T) {
	t.Run("must return `No Value`", func(t *testing.T) {
		hs := HandScore{
			Score: 900,
			Suit:  "Spades",
		}
		if hs.ToString() != "Royal Flush of Spades" {
			t.Errorf("want %v got %v", "`No Value`", hs.ToString())
		}
	})
}

func TestFindStraight(t *testing.T) {
	t.Run("must return 0 when there is no straight", func(t *testing.T) {
		cards := []int{5, 6, 7, 8, 10}
		zero := FindStraight(cards)
		if zero != 0 {
			t.Errorf("want %v got %v", 0, zero)
		}
	})

	t.Run("must return highest value of a straight", func(t *testing.T) {
		cards := []int{5, 6, 7, 8, 9}
		highest := FindStraight(cards)
		if highest != 9 {
			t.Errorf("want %v got %v", 9, highest)
		}
	})
}

func TestOrderSliceOfInt(t *testing.T) {
	t.Run("must return ordered slice of int", func(t *testing.T) {
		cards := []int{9, 8, 7, 4, 2}
		ordered := OrderSliceOfInt(cards)
		for i := 0; i < len(ordered)-1; i++ {
			if ordered[i] > ordered[i+1] {
				t.Errorf("want 2,4,7,8,9 got %v", ordered)
			}
		}
	})
}

func TestOrderSliceDescOfInt(t *testing.T) {
	t.Run("must return ordered slice of int", func(t *testing.T) {
		cards := []int{2, 4, 7, 8, 9}
		ordered := OrderSliceOfIntDesc(cards)
		for i := 0; i < len(ordered)-1; i++ {
			if ordered[i] < ordered[i+1] {
				t.Errorf("want 9,8,7,4,2 got %v", ordered)
			}
		}
	})
}

func TestDeal(t *testing.T) {
	t.Run("must deal 2 cards", func(t *testing.T) {
		cards := NewDeck()
		hand := Deal(&cards, 2)
		if len(hand) != 2 {
			t.Errorf("want 2 cards got %v", len(hand))
		}
	})
}
