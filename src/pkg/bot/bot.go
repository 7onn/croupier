package bot

import (
	"croupier/pkg/cards"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"time"
)

type InitialProbabilityMap struct {
	Card1  int
	Card2  int
	Suited bool
}

type InitialProbabiltyResult struct {
	Success int
	Try     int
}

func InitialProbabilty() error {
	probs := map[InitialProbabilityMap]InitialProbabiltyResult{}
	for i := 2; i < 15; i++ {
		for j := 2; j < i+1; j++ {
			initSuited := InitialProbabilityMap{
				Card1:  i,
				Card2:  j,
				Suited: true,
			}
			initNotSuited := InitialProbabilityMap{
				Card1:  i,
				Card2:  j,
				Suited: false,
			}
			probs[initSuited] = InitialProbabiltyResult{}
			probs[initNotSuited] = InitialProbabiltyResult{}
		}
	}

	for count := 0; count < 100000000; count++ {
		d := NewDeckForTest(count)
		if count%1000000 == 0 {
			fmt.Println("Done:", count/1000000, "%")
		}
		hand1 := cards.Deal(&d, 2)
		if hand1[0].Rank < hand1[1].Rank {
			hand1[0], hand1[1] = hand1[1], hand1[0]
		}

		hand2 := cards.Deal(&d, 2)
		if hand2[0].Rank < hand2[1].Rank {
			hand2[0], hand2[1] = hand2[1], hand2[0]
		}
		sharedCards := cards.Deal(&d, 5)
		hs1 := cards.CalculateFiveBestCards(append(sharedCards, hand1...))
		hs2 := cards.CalculateFiveBestCards(append(sharedCards, hand2...))
		i := cards.BetterThan(hs1, hs2)

		//if scores not equivalent
		if i != 0 {
			initialProbability1 := InitialProbabilityMap{
				Card1:  hand1[0].Rank,
				Card2:  hand1[1].Rank,
				Suited: hand1[0].Suit == hand1[1].Suit,
			}
			initialProbability2 := InitialProbabilityMap{
				Card1:  hand2[0].Rank,
				Card2:  hand2[1].Rank,
				Suited: hand2[0].Suit == hand2[1].Suit,
			}
			res1 := probs[initialProbability1]
			res1.Try++
			res2 := probs[initialProbability2]
			res2.Try++
			if i == -1 {
				res1.Success++
			} else {
				res2.Success++
			}
			probs[initialProbability1] = res1
			probs[initialProbability2] = res2
		}
	}
	resultString := "card1\tcard2\tsuited\tsuccess\ttry\n"
	for key, value := range probs {
		resultString += strconv.Itoa(key.Card1) + "\t" + strconv.Itoa(key.Card2) + "\t" + strconv.FormatBool(key.Suited) + "\t" + strconv.Itoa(value.Success) + "\t" + strconv.Itoa(value.Try) + "\n"
	}

	return ioutil.WriteFile("probability.csv", []byte(resultString), 0666)
}

func NewDeckForTest(i int) []cards.Card {
	suitValues := []string{"Spades", "Clubs", "Hearts", "Diamonds"}
	var d []cards.Card
	for i := 2; i < 15; i++ {
		for _, suit := range suitValues {
			d = append(d, cards.AddCard(i, suit))
		}
	}
	ShuffleForTest(d, i)
	return d
}

func ShuffleForTest(d []cards.Card, i int) {
	source := rand.NewSource(time.Now().UnixNano() + int64(i))
	r := rand.New(source)
	for i := range d {
		newPosition := r.Intn(len(d) - 1)
		d[i], d[newPosition] = d[newPosition], d[i]
	}

}
