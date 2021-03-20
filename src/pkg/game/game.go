package game

import (
	"bufio"
	"croupier/pkg/cards"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	d = cards.Shuffle(d)
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

type Coefficients struct {
	Normalisation int
	Check         float64
	Raise         float64
	Fold1         float64
	Fold2         float64
	Call          float64
	AllIn         float64
}

func (g Game) BetBlinds(botMode bool) {
	g.Rounds[len(g.Rounds)-1].Bet(g.SmallBlindTurn, g.SmallBlindValue, botMode)
	if !botMode {
		fmt.Println("\t", g.Rounds[len(g.Rounds)-1].Players[g.SmallBlindTurn].Name, "is the small blind and bet", g.SmallBlindValue)
	}

	if len(g.Rounds[len(g.Rounds)-1].Players) == 2 && g.Rounds[len(g.Rounds)-1].PlayersAllIn == 1 {
		g.Rounds[len(g.Rounds)-1].Bet(g.BigBlindTurn, g.Rounds[len(g.Rounds)-1].MaxBet, botMode)
	} else {
		g.Rounds[len(g.Rounds)-1].Bet(g.BigBlindTurn, g.BigBlindValue, botMode)
	}
	if !botMode {
		fmt.Print("\t", g.Rounds[len(g.Rounds)-1].Players[g.BigBlindTurn].Name, " is the big blind and bet ", g.BigBlindValue, "\n\n")
	}

}
func (pointerToRound *Round) Bet(playerIndex int, amount int, botMode bool) {
	bet := cards.MinInt(amount, (*pointerToRound).LimitBet-(*pointerToRound).Players[playerIndex].RoundBet)
	if (*pointerToRound).Players[playerIndex].InitialStack-(*pointerToRound).Players[playerIndex].RoundBet <= bet {
		(*pointerToRound).Pot += (*pointerToRound).Players[playerIndex].InitialStack - (*pointerToRound).Players[playerIndex].RoundBet
		(*pointerToRound).Players[playerIndex].RoundBet = (*pointerToRound).Players[playerIndex].InitialStack
		(*pointerToRound).MaxBet = cards.MaxInt((*pointerToRound).Players[playerIndex].InitialStack, (*pointerToRound).MaxBet)
		(*pointerToRound).Players[playerIndex].IsAllIn = true
		(*pointerToRound).PlayersAllIn++
		if !botMode {
			fmt.Println("\t", (*pointerToRound).Players[playerIndex].Name, " is all in!")
		}

	} else {
		(*pointerToRound).Players[playerIndex].RoundBet += bet
		(*pointerToRound).Pot += bet
		(*pointerToRound).MaxBet = (*pointerToRound).Players[playerIndex].RoundBet
	}
}

func (g *Game) EndRound(botMode bool) {
	r := (*g).Rounds[len((*g).Rounds)-1]
	// actualize the stacks

	results := r.GetResults()
	alreadyWon := 0
	for r.Pot > alreadyWon {
		for _, sp := range results {
			for pr := range sp {
				for _, roundPlayer := range r.Players {
					sp[pr].MaxWin += cards.MinInt(r.Players[sp[pr].PlayerIndex].RoundBet, roundPlayer.RoundBet)
				}
				sp[pr].MaxWin -= alreadyWon
			}
			sp = OrderResultsWithMaxWins(sp)
			r.Players[sp[0].PlayerIndex].Won = sp[0].MaxWin / len(sp)
			alreadyWon += r.Players[sp[0].PlayerIndex].Won
			if len(sp) > 1 {
				for i := 1; i < len(sp); i++ {
					r.Players[sp[i].PlayerIndex].Won = r.Players[sp[i-1].PlayerIndex].Won + (sp[i].MaxWin-sp[i-1].MaxWin)/(len(sp)-i)
					alreadyWon += r.Players[sp[i].PlayerIndex].Won
				}
			}
		}
	}
	for i := range (*g).Players {
		if (*g).Players[i].Stack > 0 {
			for playerIndex, playerInRound := range r.Players {
				if (*g).Players[i].Name == playerInRound.Name {
					(*g).Players[i].Stack -= playerInRound.RoundBet
					(*g).Players[i].Stack += playerInRound.Won
					if !botMode {
						for _, sp := range results {
							for pr := range sp {
								if sp[pr].PlayerIndex == playerIndex {
									if r.PlayersAlive > 1 {
										fmt.Print("\n\t"+playerInRound.Name, " has ")
										cards.PrintDeck(playerInRound.Cards)
										if playerInRound.Won != 0 {
											fmt.Println("\n\t"+playerInRound.Name, "won", playerInRound.Won, "with "+sp[pr].HandScore.ToString())
										}

									} else {
										fmt.Println("\n\t"+playerInRound.Name, "won", playerInRound.Won)
									}
									break
								}
							}
						}
					}

				}
			}

		}
	}

	/*for _, p := range r.players {
		if !p.hasFolded {
			fmt.Println(p.name, "has", p.cards[0].toString(), " ", p.cards[1].toString())
		}
	}*/
	// actualize the blind turns
	playersNumber := (*g).GetNumberOfPlayersAlive()
	(*g).SmallBlindTurn = ((*g).SmallBlindTurn + 1) % playersNumber
	(*g).BigBlindTurn = ((*g).BigBlindTurn + 1) % playersNumber
	(*g).DealerTurn = ((*g).DealerTurn + 1) % playersNumber
	// actualize the blind values
	if len((*g).Rounds)%10 == 0 {
		(*g).BigBlindValue *= 2
		(*g).SmallBlindValue *= 2
		if !botMode {
			fmt.Println("the blinds increases!")
			fmt.Println("Small blind: ", (*g).SmallBlindValue)
			fmt.Println("Big blind: ", (*g).BigBlindValue)
		}

	}

}

func OrderResultsWithMaxWins(prs []PlayerResult) []PlayerResult {
	for i := range prs {
		j := i
		for j >= 0 && j < len(prs)-1 && prs[j+1].MaxWin < prs[j].MaxWin {
			prs[j+1], prs[j] = prs[j], prs[j+1]
			j--
		}

	}
	return prs
}

func (r Round) GetResults() RoundResults {
	if r.PlayersAlive == 1 {
		for i := range r.Players {
			if !r.Players[i].HasFolded {
				winner := PlayerResult{
					PlayerIndex: i,
					MaxWin:      0,
				}
				results := RoundResults{}
				return append(results, []PlayerResult{winner})
			}
		}
	}
	sh := []cards.HandScore{}
	for i, p := range r.Players {
		if !p.HasFolded {
			sh = append(sh, cards.CalculateFiveBestCards(append(r.SharedCards, p.Cards...)))
			sh[len(sh)-1].PlayerIndex = i
		}
	}

	for i := range sh {
		j := i
		for j >= 0 && j < len(sh)-1 && sh[j+1].Score > sh[j].Score {
			sh[j+1], sh[j] = sh[j], sh[j+1]
			j--
		}
	}
	noswap := false
	for !noswap {
		noswap = true
		j := 0
		for j >= 0 && j < len(sh)-1 {
			x := cards.BetterThan(sh[j], sh[j+1])
			if x == 1 {
				sh[j+1], sh[j] = sh[j], sh[j+1]
				noswap = false
				sh[j+1].IsPrecedentEqual = false
				sh[j].IsPrecedentEqual = false
			}
			if x == 0 {
				sh[j+1].IsPrecedentEqual = true
			}
			j++
		}
	}
	results := RoundResults{}
	for _, h := range sh {
		playerR := PlayerResult{
			PlayerIndex: h.PlayerIndex,
			MaxWin:      0,
			HandScore:   h,
		}
		if !h.IsPrecedentEqual {
			slicePlayerResult := []PlayerResult{playerR}
			results = append(results, slicePlayerResult)
		} else {
			results[len(results)-1] = append(results[len(results)-1], playerR)
		}
	}

	return results

}

func (r Round) IsBetTurnOver(playerIndex int) bool {
	p := r.Players[playerIndex]
	if (!p.HasFolded && p.HasSpoken && (p.RoundBet == r.MaxBet || p.IsAllIn)) || r.PlayersAlive == 1 || (r.PlayersAlive-r.PlayersAllIn <= 1 && (p.RoundBet == r.MaxBet || p.IsAllIn)) {
		return true
	}

	return false
}

func ReadFromTerminal() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	text = strings.Replace(text, "\r", "", -1)

	return text
}

func LaunchGame() Game {

	sp := []Person{}
	for {
		fmt.Print("-> How many real players? (0 to 9)\n")

		text := ReadFromTerminal()

		number, err := strconv.Atoi(text)
		if err == nil {
			if number >= 0 && number < 10 {
				for i := 0; i < number; i++ {
					p := Person{Stack: 1000}
					fmt.Print("-> Name of player ", i+1, "?\n")
					name := ReadFromTerminal()
					p.Name = name
					sp = append(sp, p)
				}

				fmt.Print("-> How many virtual players? (" + strconv.Itoa(9-number) + " max)\n")
				text2 := ReadFromTerminal()

				numberOfBots, err2 := strconv.Atoi(text2)
				if err2 == nil && numberOfBots >= 0 && numberOfBots+number < 10 {
					for i := 0; i < numberOfBots; i++ {
						p := Person{Stack: 1000, IsBot: true}
						fmt.Print("-> Name of the virtual player ", i+1, "?\n")
						name := ReadFromTerminal()
						p.Name = name
						sp = append(sp, p)
					}
				} else {
					fmt.Println("Unvalid number")
					os.Exit(1)
				}
			} else {
				fmt.Println("Sorry, Only 1 to 9 players can play")
				os.Exit(1)
			}
		} else {
			fmt.Println("Unvalid number")
			os.Exit(1)
		}

		return NewGame(sp)
	}
}

func NewGame(p []Person) Game {
	g := Game{
		SmallBlindTurn:  0,
		BigBlindTurn:    1,
		DealerTurn:      len(p) - 1,
		SmallBlindValue: 10,
		BigBlindValue:   20,
		Rounds:          []Round{},
		Players:         p,
		TotalStack:      p[0].Stack * len(p),
	}

	return g
}

func (g Game) IsGameOn() bool {
	n := 0
	for _, player := range g.Players {
		if player.Stack > 0 {
			n++
		}
		if n > 1 {
			return true
		}
	}
	return false
}

func (g Game) GetNumberOfPlayersAlive() int {
	n := 0
	for _, player := range g.Players {
		if player.Stack > 0 {
			n++
		}
	}
	return n
}

func (g Game) GetWinner() string {
	for _, player := range g.Players {
		if player.Stack > 0 {
			return player.Name
		}
	}
	return ""
}

func (g Game) GetWinnerIndex() int {
	for i, player := range g.Players {
		if player.Stack > 0 {
			return i
		}
	}
	return 0
}
