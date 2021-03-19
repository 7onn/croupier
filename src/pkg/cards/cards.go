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

func BetterThan(h1 HandScore, h2 HandScore) bool {
	if h1.Score > h2.Score {
		return true
	}
	if h1.Score < h2.Score {
		return false
	}
	if h1.Card1 > h2.Card1 {
		return true
	}
	if h1.Card1 < h2.Card1 {
		return false
	}
	if h1.Card2 > h2.Card2 {
		return true
	}
	if h1.Card2 < h2.Card2 {
		return false
	}
	for i := len(h1.RemainingCards) - 1; i >= 0; i-- {
		if h1.RemainingCards[i] > h2.RemainingCards[i] {
			return true
		}
		if h1.RemainingCards[i] < h2.RemainingCards[i] {
			return false
		}
	}

	return false
}

func CalculateFiveBestCards(cardsInput []Card) HandScore {
	deck := CopyDeck(cardsInput)
	hS := HandScore{Score: -100}
	var flush []int
	var pairs []int
	var threeOfAKind []int
	var remainingCards []int
	groupOfSameColours := map[string]int{}
	groupByValues := map[int]int{}

	deck = SortCards(deck)

	for _, card := range deck {
		if groupByValues[card.Rank] == 0 {
			groupByValues[card.Rank] = 1
		} else {
			groupByValues[card.Rank]++
		}
		if groupOfSameColours[card.Suit] == 0 {
			groupOfSameColours[card.Suit] = 1
		} else {
			groupOfSameColours[card.Suit]++
		}
	}

	// find royal flush, straight flush or flush
	for suit, amount := range groupOfSameColours {
		if amount > 4 {
			hS.Suit = suit
			sc := []int{}
			for _, card := range deck {
				if card.Suit == suit {
					sc = append(sc, card.Rank)
				}
			}

			//if royal flush
			if sc[len(sc)-5] == 10 {
				hS.Score = 900
			} else {
				//if straight
				if highestCard := FindStraight(sc); highestCard != 0 {
					hS.Score = 800
					hS.Card1 = highestCard
				} else {
					flush = sc[len(sc)-5:]
				}
			}
		}
	}

	// find four of a kind
	if hS.Score == -100 {
		remains := []int{}
		for cardRank, amount := range groupByValues {
			if amount == 4 {
				hS.Score = 700
				hS.Card1 = cardRank
			} else {
				remains = append(remains, cardRank)
			}
		}
		if hS.Score == 700 {
			hS.RemainingCards = remains[len(remains)-1:]
		}
	}

	// find full house - three of a kind - 1 pair - 2 pairs
	if hS.Score == -100 {
		for cardRank, amount := range groupByValues {
			if amount == 3 {
				threeOfAKind = append(threeOfAKind, cardRank)
			} else {
				if amount == 2 {
					pairs = append(pairs, cardRank)
				} else {
					remainingCards = append(remainingCards, cardRank)
					remainingCards = OrderSliceOfInt(remainingCards)
				}
			}
		}
		threeOfAKind = OrderSliceOfInt(threeOfAKind)

		if len(threeOfAKind) == 2 {
			pairs = append(pairs, threeOfAKind[0])
			threeOfAKind = threeOfAKind[1:]
		}

		pairs = OrderSliceOfInt(pairs)

		if len(threeOfAKind) == 1 {
			if len(pairs) > 0 {
				hS.Score = 600
				hS.Card1 = threeOfAKind[0]
				hS.Card2 = pairs[len(pairs)-1]
			}
		}

	}

	if hS.Score == -100 && len(flush) > 0 {
		hS.Score = 500
		hS.RemainingCards = flush
		hS.Card1 = flush[len(flush)-1]
	}

	// find straight including the straight Ace 2 3 4 5 (Ace=14)
	if hS.Score == -100 {
		cardValues := []int{}
		for _, card := range deck {
			cardValues = append(cardValues, card.Rank)
		}
		if highestCard := FindStraight(cardValues); highestCard != 0 {
			hS.Score = 400
			hS.Card1 = highestCard
		}

	}

	// remains three of a kind, single pair, two pairs, or high card
	if hS.Score == -100 {
		if len(threeOfAKind) > 0 {
			hS.Score = 300
			hS.Card1 = threeOfAKind[0]
			hS.RemainingCards = remainingCards[len(remainingCards)-2:]

		} else {
			if len(pairs) > 2 {
				remainingCards = append(remainingCards, pairs[(len(pairs)-3)])
				remainingCards = OrderSliceOfInt(remainingCards)
				pairs = pairs[len(pairs)-2:]
			}
			remainingCards = OrderSliceOfInt(remainingCards)
			switch len(pairs) {
			case 2:
				hS.Score = 200
				hS.Card1 = pairs[len(pairs)-1]
				hS.Card2 = pairs[len(pairs)-2]
				hS.RemainingCards = remainingCards[len(remainingCards)-1:]
			case 1:
				hS.Score = 100
				hS.Card1 = pairs[0]
				hS.RemainingCards = remainingCards[len(remainingCards)-3:]
			case 0:
				hS.Score = 0
				hS.Card1 = remainingCards[len(remainingCards)-1]
				hS.RemainingCards = remainingCards[len(remainingCards)-5 : len(remainingCards)-1]
			}

		}
	}

	return hS
}
