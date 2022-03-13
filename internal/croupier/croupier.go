package croupier

import (
	"fmt"
	"time"

	ws "github.com/7onn/croupier/pkg/websocket"
)

// Croupier exposes all functionalities of the Croupier service.
type Croupier interface {
	Ping() bool
	PlayPoker(client *ws.Client)
}

// Broker manages the internal state of the Croupier service.
type Broker struct{}

// New initializes a new Croupier service.
func New() *Broker {
	r := &Broker{}
	return r
}

// Ping checks to see if the croupier's database is responding.
func (brk *Broker) Ping() bool {
	// This function would check the croupier's dependencies (datastores and whatnot); useful for Kubernetes probes
	return true
}

func (brk *Broker) PlayPoker(client *ws.Client) {
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
