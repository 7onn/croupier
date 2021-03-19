package cards

import (
	"math/rand"
	"time"
)

type Card struct {
	Suit string
	Rank int
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
	return suitMapping[c.Suit] + valuesMapping[c.Rank]
}

type HandScore struct {
	Score            int
	Suit             string
	Card1            int
	Card2            int
	RemainingCards   []int
	PlayerIndex      int
	IsPrecedentEqual bool
}

func (s HandScore) ToString() string {
	valuesMappingInt := map[int]string{
		2:  "Two",
		3:  "Three",
		4:  "Four",
		5:  "Five",
		6:  "Six",
		7:  "Seven",
		8:  "Eight",
		9:  "Nine",
		10: "Ten",
		11: "Jalet",
		12: "Queen",
		13: "King",
		14: "Ace",
	}
	valuesMapping := map[int]string{
		0:    "a High Card " + valuesMappingInt[s.Card1],
		100:  "a Single Pair of " + valuesMappingInt[s.Card1] + "s",
		200:  "Two Pairs of " + valuesMappingInt[s.Card1] + "s" + " and " + valuesMappingInt[s.Card2] + "s",
		300:  "a Three Of A Kind of " + valuesMappingInt[s.Card1] + "s",
		400:  "a " + valuesMappingInt[s.Card1] + "-high Straigth",
		500:  "a Flush of " + s.Suit,
		600:  "a Full House, " + valuesMappingInt[s.Card1] + "s" + " over " + valuesMappingInt[s.Card2] + "s",
		700:  "a Four Of A Kind of " + valuesMappingInt[s.Card1] + "s",
		800:  "a " + valuesMappingInt[s.Card1] + "-high Straight Flush of " + s.Suit,
		900:  "Royal Flush of " + s.Suit,
		-100: "No Value",
	}

	return valuesMapping[s.Score]

}

func AddCard(rank int, suit string) Card {
	return Card{
		Rank: rank,
		Suit: suit,
	}
}

func NewDeck() []Card {
	suits := []string{"Spades", "Clubs", "Hearts", "Diamonds"}
	var deck []Card
	for _, suit := range suits {
		for rank := 2; rank <= 14; rank++ {
			deck = append(deck, AddCard(rank, suit))
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

func FindStraight(c []int) int {
	best := []int{}
	curr := []int{c[0]}
	for i := range c {
		if i > 0 && c[i] == c[i-1]+1 {
			curr = append(curr, c[i])
		} else {
			if i > 0 && c[i] != c[i-1] {
				if len(curr) == 5 || (len(curr) == 4 && curr[0] == 2 && c[len(c)-1] == 14) {
					best = curr
				}
				curr = []int{c[i]}
			}
		}
	}
	if len(curr) >= 5 || (len(curr) == 4 && curr[0] == 2 && c[len(c)-1] == 14) {
		best = curr
	}
	if len(best) > 0 {
		return best[len(best)-1]
	}
	return 0
}

func CopyDeck(cards []Card) []Card {
	copy := []Card{}
	for _, card := range cards {
		copy = append(copy, card)
	}

	return copy
}

func SortCards(c []Card) []Card {
	for i := range c {
		j := i
		for j >= 0 && j < len(c)-1 && c[j+1].Rank < c[j].Rank {
			c[j+1], c[j] = c[j], c[j+1]
			j--
		}

	}
	return c
}

func OrderSliceOfInt(s []int) []int {
	for i := range s {
		j := i
		for j >= 0 && j < len(s)-1 && s[j+1] < s[j] {
			s[j+1], s[j] = s[j], s[j+1]
			j--
		}

	}
	return s
}

func OrderSliceOfIntDesc(s []int) []int {
	for i := range s {
		j := i
		for j >= 0 && j < len(s)-1 && s[j+1] > s[j] {
			s[j+1], s[j] = s[j], s[j+1]
			j--
		}

	}
	return s
}

func Deal(d *[]Card, handSize int) []Card {
	hand := (*d)[:handSize]
	(*d) = (*d)[handSize:]
	return hand
}
