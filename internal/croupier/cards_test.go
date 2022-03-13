package croupier

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

func TestBetterThan(t *testing.T) {
	t.Run("royal flush > straight flush", func(t *testing.T) {
		rf := []Card{
			{
				Suit: "Spades",
				Rank: 10,
			},
			{
				Suit: "Spades",
				Rank: 11,
			},
			{
				Suit: "Spades",
				Rank: 12,
			},
			{
				Suit: "Spades",
				Rank: 13,
			},
			{
				Suit: "Spades",
				Rank: 14,
			},
		}

		sf := []Card{
			{
				Suit: "Spades",
				Rank: 6,
			},
			{
				Suit: "Spades",
				Rank: 7,
			},
			{
				Suit: "Spades",
				Rank: 8,
			},
			{
				Suit: "Spades",
				Rank: 9,
			},
			{
				Suit: "Spades",
				Rank: 10,
			},
		}

		hs1 := CalculateFiveBestCards(rf)
		hs2 := CalculateFiveBestCards(sf)

		if BetterThan(hs1, hs2) != -1 {
			t.Errorf("royal flush is better than straight flush")
		}

	})

	t.Run("straight flush > four of a kind", func(t *testing.T) {
		sf := []Card{
			{
				Suit: "Spades",
				Rank: 6,
			},
			{
				Suit: "Spades",
				Rank: 7,
			},
			{
				Suit: "Spades",
				Rank: 8,
			},
			{
				Suit: "Spades",
				Rank: 9,
			},
			{
				Suit: "Spades",
				Rank: 10,
			},
		}

		fok := []Card{
			{
				Suit: "Spades",
				Rank: 14,
			},
			{
				Suit: "Diamonds",
				Rank: 14,
			},
			{
				Suit: "Hearts",
				Rank: 14,
			},
			{
				Suit: "Spades",
				Rank: 14,
			},
			{
				Suit: "Spades",
				Rank: 10,
			},
		}

		hs1 := CalculateFiveBestCards(sf)
		hs2 := CalculateFiveBestCards(fok)

		if BetterThan(hs1, hs2) != -1 {
			t.Errorf("straight flush should beat four of a kind")
		}

	})
}

func TestCalculateFiveBestCards(t *testing.T) {
	t.Run("must return nothing", func(t *testing.T) {
		cards := []Card{
			{
				Suit: "Hearts",
				Rank: 2,
			},
			{
				Suit: "Diamonds",
				Rank: 3,
			},
			{
				Suit: "Spades",
				Rank: 5,
			},
			{
				Suit: "Cubs",
				Rank: 6,
			},
			{
				Suit: "Spades",
				Rank: 8,
			},
		}
		hs := CalculateFiveBestCards(cards)
		if hs.Score != 0 {
			t.Errorf("wanted score %v got %v", 0, hs.Score)
		}
	})

	t.Run("must return royal flush", func(t *testing.T) {
		cards := []Card{
			{
				Suit: "Spades",
				Rank: 10,
			},
			{
				Suit: "Spades",
				Rank: 11,
			},
			{
				Suit: "Spades",
				Rank: 12,
			},
			{
				Suit: "Spades",
				Rank: 13,
			},
			{
				Suit: "Spades",
				Rank: 14,
			},
		}
		hs := CalculateFiveBestCards(cards)
		if hs.Score != 900 {
			t.Errorf("wanted score %v got %v", 900, hs.Score)
		}
	})

	t.Run("must return straight flush", func(t *testing.T) {
		cards := []Card{
			{
				Suit: "Spades",
				Rank: 6,
			},
			{
				Suit: "Spades",
				Rank: 7,
			},
			{
				Suit: "Spades",
				Rank: 8,
			},
			{
				Suit: "Spades",
				Rank: 9,
			},
			{
				Suit: "Spades",
				Rank: 10,
			},
		}
		hs := CalculateFiveBestCards(cards)
		if hs.Score != 800 {
			t.Errorf("wanted score %v got %v", 800, hs.Score)
		}
	})

	t.Run("must return four of a kind", func(t *testing.T) {
		cards := []Card{
			{
				Suit: "Hearts",
				Rank: 2,
			},
			{
				Suit: "Diamonds",
				Rank: 2,
			},
			{
				Suit: "Spades",
				Rank: 2,
			},
			{
				Suit: "Cubs",
				Rank: 2,
			},
			{
				Suit: "Spades",
				Rank: 8,
			},
		}
		hs := CalculateFiveBestCards(cards)
		if hs.Score != 700 {
			t.Errorf("wanted score %v got %v", 700, hs.Score)
		}
	})

	t.Run("must return full house", func(t *testing.T) {
		cards := []Card{
			{
				Suit: "Hearts",
				Rank: 2,
			},
			{
				Suit: "Diamonds",
				Rank: 2,
			},
			{
				Suit: "Spades",
				Rank: 2,
			},
			{
				Suit: "Cubs",
				Rank: 3,
			},
			{
				Suit: "Spades",
				Rank: 3,
			},
		}
		hs := CalculateFiveBestCards(cards)
		if hs.Score != 600 {
			t.Errorf("wanted score %v got %v", 600, hs.Score)
		}
	})

	t.Run("must return straight", func(t *testing.T) {
		cards := []Card{
			{
				Suit: "Spades",
				Rank: 14,
			},
			{
				Suit: "Cubs",
				Rank: 2,
			},
			{
				Suit: "Hearts",
				Rank: 3,
			},
			{
				Suit: "Diamonds",
				Rank: 4,
			},
			{
				Suit: "Spades",
				Rank: 5,
			},
		}
		hs := CalculateFiveBestCards(cards)
		if hs.Score != 400 {
			t.Errorf("wanted score %v got %v", 400, hs.Score)
		}
	})

	t.Run("must return flush", func(t *testing.T) {
		cards := []Card{
			{
				Suit: "Spades",
				Rank: 2,
			},
			{
				Suit: "Spades",
				Rank: 3,
			},
			{
				Suit: "Spades",
				Rank: 5,
			},
			{
				Suit: "Spades",
				Rank: 6,
			},
			{
				Suit: "Spades",
				Rank: 8,
			},
		}
		hs := CalculateFiveBestCards(cards)
		if hs.Score != 500 {
			t.Errorf("wanted score %v got %v", 500, hs.Score)
		}
	})

	t.Run("must return three of a kind", func(t *testing.T) {
		cards := []Card{
			{
				Suit: "Hearts",
				Rank: 2,
			},
			{
				Suit: "Diamonds",
				Rank: 2,
			},
			{
				Suit: "Spades",
				Rank: 2,
			},
			{
				Suit: "Cubs",
				Rank: 7,
			},
			{
				Suit: "Spades",
				Rank: 8,
			},
		}
		hs := CalculateFiveBestCards(cards)
		if hs.Score != 300 {
			t.Errorf("wanted score %v got %v", 300, hs.Score)
		}
	})

	t.Run("must return two pairs", func(t *testing.T) {
		cards := []Card{
			{
				Suit: "Spades",
				Rank: 3,
			},
			{
				Suit: "Spades",
				Rank: 12,
			},
			{
				Suit: "Diamonds",
				Rank: 3,
			},
			{
				Suit: "Hearts",
				Rank: 10,
			},
			{
				Suit: "Hearts",
				Rank: 12,
			},
		}
		hs := CalculateFiveBestCards(cards)
		if hs.Score != 200 {
			t.Errorf("wanted score %v got %v", 200, hs.Score)
		}
	})

	t.Run("must return a pair", func(t *testing.T) {
		cards := []Card{
			{
				Suit: "Hearts",
				Rank: 2,
			},
			{
				Suit: "Diamonds",
				Rank: 2,
			},
			{
				Suit: "Spades",
				Rank: 3,
			},
			{
				Suit: "Cubs",
				Rank: 7,
			},
			{
				Suit: "Spades",
				Rank: 8,
			},
		}
		hs := CalculateFiveBestCards(cards)
		if hs.Score != 100 {
			t.Errorf("wanted score %v got %v", 100, hs.Score)
		}
	})

}
