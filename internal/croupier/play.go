package croupier

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	ws "github.com/7onn/croupier/pkg/websocket"
)

func ExpoFunction(odds float64, index float64) float64 {
	return math.Exp(index * (odds - 1.0))
}

func RaiseExpo(x float64) float64 {
	return (1.5 - 0.5*x) * math.Exp(x-1.0)
}

func NormaleLawRandomization(risk float64, offset float64) float64 {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	riskInt := int(math.Ceil(risk * 100))
	randNumber := r.Intn(riskInt*2) - riskInt
	x := float64(randNumber)/100 + offset
	if x < 0 {
		return math.Exp(-0.5 * x * x)
	}

	return 2.0 - math.Exp(-0.5*x*x)

}

func RatioStacks(stack float64, opponentStack float64) float64 {
	return (math.Log10(stack/opponentStack) + 2.0) / 2.0
}

func RatioToWintoBet(amountToBet float64, pot float64) float64 {
	return pot / amountToBet
}

func PlayPoker(client *ws.Client) {
	g := LaunchGame()

	coeff := Coefficients{
		Check:         0.65,
		Raise:         0.3,
		Call:          0.3,
		Normalisation: 3,
		Fold1:         0.5,
		Fold2:         1.5,
		AllIn:         0.95,
	}

	for g.IsGameOn() {

		g.NewRound(client, false)
		ShowCards(client, g.Rounds[len(g.Rounds)-1].Players)
		g.BetBlinds(client, false)

		numberOfPlayers := g.GetNumberOfPlayersAlive()
		turnToPlay := (g.SmallBlindTurn + 2) % numberOfPlayers
		for !g.Rounds[len(g.Rounds)-1].IsBetTurnOver(turnToPlay) {
			Play(client, &(g.Rounds[len(g.Rounds)-1]), turnToPlay, g.BigBlindValue, float64(g.TotalStack), coeff, false)
			turnToPlay = (turnToPlay + 1) % numberOfPlayers
		}
		for g.Rounds[len(g.Rounds)-1].PlayersAlive > 1 && len(g.Rounds[len(g.Rounds)-1].SharedCards) < 5 {
			switch len(g.Rounds[len(g.Rounds)-1].SharedCards) {
			case 0:
				client.Send <- []byte("\n\tHere comes the flop...\n\n")
				g.Rounds[len(g.Rounds)-1].SharedCards = Deal(&g.Rounds[len(g.Rounds)-1].Deck, 3)

				time.Sleep(time.Second * 1)
				client.Send <- []byte(g.Rounds[len(g.Rounds)-1].SharedCards[0].ToString())

				time.Sleep(time.Second * 1)
				client.Send <- []byte(g.Rounds[len(g.Rounds)-1].SharedCards[1].ToString())

				time.Sleep(time.Second * 1)
				client.Send <- []byte(g.Rounds[len(g.Rounds)-1].SharedCards[2].ToString())

				time.Sleep(time.Second * 1)
			case 3:
				client.Send <- []byte("Here comes the turn...")
				g.Rounds[len(g.Rounds)-1].SharedCards = append(g.Rounds[len(g.Rounds)-1].SharedCards, Deal(&g.Rounds[len(g.Rounds)-1].Deck, 1)...)
				time.Sleep(time.Second * 1)
				client.Send <- []byte(g.Rounds[len(g.Rounds)-1].SharedCards[3].ToString())
				time.Sleep(time.Second * 1)

				client.Send <- []byte("CARDS:")
				PrintDeck(client, g.Rounds[len(g.Rounds)-1].SharedCards)
			case 4:
				client.Send <- []byte("Here comes the river...")
				g.Rounds[len(g.Rounds)-1].SharedCards = append(g.Rounds[len(g.Rounds)-1].SharedCards, Deal(&g.Rounds[len(g.Rounds)-1].Deck, 1)...)
				time.Sleep(time.Second * 1)
				PrintDeck(client, g.Rounds[len(g.Rounds)-1].SharedCards)
			}

			for i := range g.Rounds[len(g.Rounds)-1].Players {
				g.Rounds[len(g.Rounds)-1].Players[i].HasSpoken = false
			}
			turnToPlay := g.SmallBlindTurn
			for !g.Rounds[len(g.Rounds)-1].IsBetTurnOver(turnToPlay) {
				Play(client, &(g.Rounds[len(g.Rounds)-1]), turnToPlay, g.BigBlindValue, float64(g.TotalStack), coeff, false)
				turnToPlay = (turnToPlay + 1) % numberOfPlayers
			}
		}
		g.EndRound(client, false)
		time.Sleep(time.Second * 3)

	}
	client.Send <- []byte(fmt.Sprintf("Congratulations %s you won", g.GetWinner()))
}

func PrintDeck(client *ws.Client, d []Card) {
	s := ""
	for _, card := range d {
		s += card.ToString() + " "
	}
	client.Send <- []byte(s)
}

func ShowCards(client *ws.Client, sp []Player) {
	realPlayers := 0
	for _, p := range sp {
		if !p.IsBot {
			realPlayers++
		}
	}
	for _, p := range sp {
		if realPlayers > 1 {
			if p.IsBot {
				client.Send <- []byte("---")
				client.Send <- []byte(fmt.Sprintf("%s, your cards are: %s %s", p.Name, p.Cards[0].ToString(), p.Cards[1].ToString()))
			} else {
				c := ""
				for strings.ToLower(c) != "ok" {
					client.Send <- []byte(fmt.Sprintf("%s, press OK to see your cards", p.Name))
				}
				client.Send <- []byte("---")
				client.Send <- []byte(fmt.Sprintf("%s, your cards are: %s %s", p.Name, p.Cards[0].ToString(), p.Cards[1].ToString()))
				c = ""
				for strings.ToLower(c) != "ok" {
					client.Send <- []byte(fmt.Sprintf("%s, press OK to hide your cards", p.Name))
				}

			}
		} else {
			if realPlayers == 1 {
				if !p.IsBot {
					client.Send <- []byte("---")
					client.Send <- []byte(fmt.Sprintf("%v, your cards are: %v %v", p.Name, p.Cards[0].ToString(), p.Cards[1].ToString()))
				}
			} else {
				client.Send <- []byte("---")
				client.Send <- []byte(fmt.Sprintf("%v, your cards are: %v %v", p.Name, p.Cards[0].ToString(), p.Cards[1].ToString()))
			}

		}
	}
}

func Play(client *ws.Client, r *Round, playerIndex int, bigBlind int, totalStacks float64, coeff Coefficients, botMode bool) {
	pp := &((*r).Players[playerIndex])
	p := (*pp)
	if p.IsAllIn && !botMode {
		fmt.Print("\t", p.Name, " is all in\n")
		client.Send <- []byte(fmt.Sprintf("%v is all in", p.Name))
	}
	if !p.IsAllIn && !p.HasFolded {
		if !p.IsBot {
			client.Send <- []byte(fmt.Sprintf("%v, What do you want to do? You have $%v", p.Name, p.InitialStack-p.RoundBet))
			if (*r).MaxBet == p.RoundBet {
				client.Send <- []byte("Check (ch)")
			} else {
				client.Send <- []byte(fmt.Sprintf("> Call %v (ca)", (r.MaxBet - p.RoundBet)))
			}
			if p.InitialStack > (*r).MaxBet {
				client.Send <- []byte("> Raise (r)")
			}
			client.Send <- []byte("> Fold (f)")
			validInput := false
			for !validInput {
				validInput = true
				_, c, _ := client.Conn.ReadMessage()
				switch string(c) {
				case "check", "ch":
					if (*r).MaxBet != p.RoundBet {
						client.Send <- []byte(fmt.Sprintf("%v checks", p.Name))
						validInput = false
					}
				case "call", "ca":
					if (*r).MaxBet != p.RoundBet {
						client.Send <- []byte(fmt.Sprintf("%v calls", p.Name))
						(*r).Bet(client, playerIndex, (*r).MaxBet-p.RoundBet, false)
					}
				case "fold", "f":
					(*pp).HasFolded = true
					(*r).PlayersAlive--
					client.Send <- []byte(fmt.Sprintf("* %v folds", p.Name))
				case "raise", "r":
					if p.InitialStack > (*r).MaxBet {
						valueRaise := 0
						for valueRaise == 0 {
							if p.InitialStack-(*r).MaxBet > bigBlind {
								client.Send <- []byte(fmt.Sprintf("---> How much do you want to raise? ( $%v, ~ $%d )", bigBlind, p.InitialStack-(*r).MaxBet))
								_, text, _ := client.Conn.ReadMessage()
								client.Send <- text
								number, err := strconv.Atoi(string(text))
								if err == nil {
									if number >= bigBlind {
										valueRaise = number
										(*r).Bet(client, playerIndex, (*r).MaxBet-p.RoundBet+number, false)
									}
								}
							} else {
								(*r).Bet(client, playerIndex, p.InitialStack-p.RoundBet, false)
							}
						}
						client.Send <- []byte(fmt.Sprintf("%s raises $%d", p.Name, valueRaise))
					}
				default:
					validInput = false
				}
				if !validInput {
					client.Send <- []byte("Invalid answer")
				}
			}
		} else {
			var decision float64

			if (*pp).HasSpoken {
				decision = (*pp).Decision
			} else {
				odds := CalculateOdds(p.Cards, r.SharedCards)
				decision = odds * NormaleLawRandomization(0.7, 0)
				(*pp).Decision = decision
			}

			if (*r).MaxBet == p.RoundBet {
				if decision < coeff.Check {
					if !botMode {
						client.Send <- []byte(fmt.Sprintf("%v checks", p.Name))
					}
				} else {
					raiseValue := MaxInt(int(math.Ceil(ExpoFunction(decision, 3.0)*(totalStacks/float64(len(r.Players)))*coeff.Raise*NormaleLawRandomization(1, 0))), bigBlind)
					raiseValue -= raiseValue % (bigBlind / 2)
					(*r).Bet(client, playerIndex, raiseValue, botMode)
					if !botMode {
						client.Send <- []byte(fmt.Sprintf("%s raises %d", p.Name, raiseValue))
					}

				}
			} else {
				toCallRatio := float64((*r).MaxBet-p.RoundBet) / (totalStacks * math.Exp(float64(p.RoundBet)/(totalStacks/float64(len(r.Players)))))

				limit1 := coeff.Fold1
				limit2 := RaiseExpo(toCallRatio * coeff.Fold2)

				if decision < limit1 || (decision < limit2) /* && (float64(p.initialStack-p.roundBet) > 0.05*totalStacks*/ {
					(*pp).HasFolded = true
					(*r).PlayersAlive--
					if !botMode {
						client.Send <- []byte(fmt.Sprintf("* %s folds", p.Name))
					}

				} else {
					overFold := decision - limit2
					if overFold < coeff.Call {
						(*r).Bet(client, playerIndex, (*r).MaxBet-p.RoundBet, botMode)
						if !botMode {
							client.Send <- []byte(fmt.Sprintf("%s calls", p.Name))
						}

					} else {
						var raiseValue int
						if decision > coeff.AllIn {
							raiseValue = p.InitialStack - p.RoundBet
						} else {
							raiseValue = MaxInt(int(math.Ceil(ExpoFunction(overFold, 3.0)*(totalStacks/float64(len(r.Players)))*coeff.Raise*NormaleLawRandomization(0.5, 0))), bigBlind)
							raiseValue -= raiseValue % (bigBlind / 2)
						}

						(*r).Bet(client, playerIndex, (*r).MaxBet-p.RoundBet+raiseValue, botMode)
						if !botMode {
							client.Send <- []byte(fmt.Sprintf("%s raises %d", p.Name, raiseValue))
						}

					}
				}
			}
		}
		(*pp).HasSpoken = true

	}

}

type Round struct {
	SharedCards  []Card
	Players      []Player
	Deck         []Card
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
	HandScore   HandScore
}

type SidePot struct {
	Value         int
	PlayerIndexes []int
}

type Player struct {
	Name         string
	Cards        []Card
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

func (g *Game) NewRound(client *ws.Client, botMode bool) {
	d := NewDeck()
	d = Shuffle(d)
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
			c := Deal(&d, 2)
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
		client.Send <- []byte("---")
		client.Send <- []byte(fmt.Sprintf("ROUND %d", len((*g).Rounds)))
		client.Send <- []byte(fmt.Sprintf("Players: %d", len((*g).Rounds)))
		for _, rp := range (*g).Rounds[len((*g).Rounds)-1].Players {
			client.Send <- []byte(fmt.Sprintf("%s stack: $%d", rp.Name, rp.InitialStack))
		}
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

func (g Game) BetBlinds(client *ws.Client, botMode bool) {
	g.Rounds[len(g.Rounds)-1].Bet(client, g.SmallBlindTurn, g.SmallBlindValue, botMode)
	if !botMode {
		client.Send <- []byte(fmt.Sprintf("%s is the small blind and bet $%v", g.Rounds[len(g.Rounds)-1].Players[g.SmallBlindTurn].Name, g.SmallBlindValue))
	}

	if len(g.Rounds[len(g.Rounds)-1].Players) == 2 && g.Rounds[len(g.Rounds)-1].PlayersAllIn == 1 {
		g.Rounds[len(g.Rounds)-1].Bet(client, g.BigBlindTurn, g.Rounds[len(g.Rounds)-1].MaxBet, botMode)
	} else {
		g.Rounds[len(g.Rounds)-1].Bet(client, g.BigBlindTurn, g.BigBlindValue, botMode)
	}
	if !botMode {
		client.Send <- []byte(fmt.Sprintf("%s is the big blind and bet $%v", g.Rounds[len(g.Rounds)-1].Players[g.SmallBlindTurn].Name, g.BigBlindValue))
	}

}

func (g *Game) EndRound(client *ws.Client, botMode bool) {
	r := (*g).Rounds[len((*g).Rounds)-1]
	// actualize the stacks

	results := r.GetResults()
	alreadyWon := 0
	for r.Pot > alreadyWon {
		for _, sp := range results {
			for pr := range sp {
				for _, roundPlayer := range r.Players {
					sp[pr].MaxWin += MinInt(r.Players[sp[pr].PlayerIndex].RoundBet, roundPlayer.RoundBet)
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
										client.Send <- []byte(fmt.Sprintf("%s has", playerInRound.Name))
										PrintDeck(client, playerInRound.Cards)
										if playerInRound.Won != 0 {
											client.Send <- []byte(fmt.Sprintf("%s won %d with %s", playerInRound.Name, playerInRound.Won, sp[pr].HandScore.ToString()))
										}

									} else {
										client.Send <- []byte(fmt.Sprintf("%s won %d", playerInRound.Name, playerInRound.Won))
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
			client.Send <- []byte("the blind increases!")
			client.Send <- []byte(fmt.Sprintf("Small blind: %d", (*g).SmallBlindValue))
			client.Send <- []byte(fmt.Sprintf("Big blind: %d", (*g).BigBlindValue))
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
	sh := []HandScore{}
	for i, p := range r.Players {
		if !p.HasFolded {
			sh = append(sh, CalculateFiveBestCards(append(r.SharedCards, p.Cards...)))
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
			x := BetterThan(sh[j], sh[j+1])
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

var bots = []string{"Fulane", "Ciclane", "Falsiane", "Simulane", "Pamela"}

func LaunchGame() Game {

	sp := []Person{}

	p := Person{Stack: 1000}

	if p.Name = os.Getenv("PLAYER_NAME"); p.Name == "" {
		p.Name = "tom"
	}

	sp = append(sp, p)

	numberOfBots, err2 := strconv.Atoi("5")
	if err2 == nil && numberOfBots >= 0 && numberOfBots+1 < 10 {
		for i := 0; i < numberOfBots; i++ {
			p := Person{Stack: 1000, IsBot: true}
			p.Name = bots[i]
			sp = append(sp, p)
		}
	} else {
		fmt.Println("Unvalid number")
		os.Exit(1)
	}

	return NewGame(sp)

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

func (r *Round) Bet(client *ws.Client, playerIndex int, amount int, botMode bool) {
	bet := MinInt(amount, (*r).LimitBet-(*r).Players[playerIndex].RoundBet)
	if (*r).Players[playerIndex].InitialStack-(*r).Players[playerIndex].RoundBet <= bet {
		(*r).Pot += (*r).Players[playerIndex].InitialStack - (*r).Players[playerIndex].RoundBet
		(*r).Players[playerIndex].RoundBet = (*r).Players[playerIndex].InitialStack
		(*r).MaxBet = MaxInt((*r).Players[playerIndex].InitialStack, (*r).MaxBet)
		(*r).Players[playerIndex].IsAllIn = true
		(*r).PlayersAllIn++
		if !botMode {
			client.Send <- []byte(fmt.Sprintf("%v is all in!", (*r).Players[playerIndex].Name))
		}

	} else {
		(*r).Players[playerIndex].RoundBet += bet
		(*r).Pot += bet
		(*r).MaxBet = (*r).Players[playerIndex].RoundBet
	}
}
