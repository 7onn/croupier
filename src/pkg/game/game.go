package game

import (
	"croupier/pkg/cards"
)

func CalculateFiveBestCards(cardsInput []cards.Card) cards.HandScore {
	deck := cards.CopyDeck(cardsInput)
	hS := cards.HandScore{Score: -100}
	var flush []int
	var pairs []int
	var threeOfAKind []int
	var remainingCards []int
	groupOfSameColours := map[string]int{}
	groupByValues := map[int]int{}

	deck = cards.SortCards(deck)

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
				if highestCard := cards.FindStraight(sc); highestCard != 0 {
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
					remainingCards = cards.OrderSliceOfInt(remainingCards)
				}
			}
		}
		threeOfAKind = cards.OrderSliceOfInt(threeOfAKind)

		if len(threeOfAKind) == 2 {
			pairs = append(pairs, threeOfAKind[0])
			threeOfAKind = threeOfAKind[1:]
		}

		pairs = cards.OrderSliceOfInt(pairs)

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
		if highestCard := cards.FindStraight(cardValues); highestCard != 0 {
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
				remainingCards = cards.OrderSliceOfInt(remainingCards)
				pairs = pairs[len(pairs)-2:]
			}
			remainingCards = cards.OrderSliceOfInt(remainingCards)
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
