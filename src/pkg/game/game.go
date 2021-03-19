package game

import (
	"croupier/pkg/cards"
	"fmt"
)

type Round struct {
	SharedCards  []cards.Card
	Players      []Player
	Deck         []cards.Card
	Pot          int
	MaxBet       int
	PlayersAlive int
	PlayersAllIn int
	LimitBet     int
}

type RoundResults = [][]PlayerResult

type Result struct {
	Position      int
	PlayerResults []PlayerResult
}

type PlayerResult struct {
	PlayerIndex int
	MaxWin      int
	HandScore   cards.HandScore
}

type SidePot struct {
	Value         int
	PlayerIndexes []int
}

type Player struct {
	Name         string
	Cards        []cards.Card
	Won          int
	RoundBet     int
	InitialStack int
	HasSpoken    bool
	HasFolded    bool
	IsAllIn      bool
	IsBot        bool
	Decision     float64
}

type Person struct {
	Name  string
	Stack int
	IsBot bool
}

type Game struct {
	Players         []Person
	Rounds          []Round
	SmallBlindTurn  int
	BigBlindTurn    int
	DealerTurn      int
	SmallBlindValue int
	BigBlindValue   int
	TotalStack      int
}

func (g *Game) NewRound(botMode bool) {
	d := cards.NewDeck()
	players := []Player{}
	biggestStack := 0
	secondStack := 0
	for _, person := range (*g).Players {
		if person.Stack > 0 {
			if biggestStack < person.Stack {
				secondStack = biggestStack
				biggestStack = person.Stack
			} else {
				if secondStack < person.Stack {
					secondStack = person.Stack
				}
			}
			c := cards.Deal(&d, 2)
			player := Player{
				Name:         person.Name,
				InitialStack: person.Stack,
				Cards:        c,
				RoundBet:     0,
				HasSpoken:    false,
				HasFolded:    false,
				IsAllIn:      false,
				IsBot:        person.IsBot,
			}

			players = append(players, player)

		}
	}

	r := Round{
		Deck:         d,
		Players:      players,
		PlayersAlive: len(players),
		LimitBet:     secondStack,
	}
	(*g).Rounds = append((*g).Rounds, r)
	if !botMode {
		fmt.Println()
		fmt.Println("ROUND", len((*g).Rounds))
		fmt.Println("\tPlayers:")
		for _, rp := range (*g).Rounds[len((*g).Rounds)-1].Players {
			fmt.Println("\t\t", rp.Name, "stack:", rp.InitialStack)
		}

		fmt.Println()
	}

}
