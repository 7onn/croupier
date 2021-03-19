package game

import (
	"croupier/pkg/cards"
	"testing"
)

func TestCalculateFiveBestCards(t *testing.T) {
	t.Run("must return nothing", func(t *testing.T) {
		cards := []cards.Card{
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
		cards := []cards.Card{
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
		cards := []cards.Card{
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
		cards := []cards.Card{
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
		cards := []cards.Card{
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
		cards := []cards.Card{
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
		cards := []cards.Card{
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
		cards := []cards.Card{
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
		cards := []cards.Card{
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
		cards := []cards.Card{
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
