package poker

// import (
// 	"bufio"
// 	"croupier/pkg/bot"
// 	"croupier/pkg/cards"
// 	"croupier/pkg/game"
// 	"croupier/pkg/odds"
// 	"fmt"
// 	"math"
// 	"math/rand"
// 	"os"
// 	"strconv"
// 	"strings"
// 	"time"
// )

// func PlayPoker() {
// 	g := game.LaunchGame()

// 	coeff := game.Coefficients{
// 		Check:         0.65,
// 		Raise:         0.3,
// 		Call:          0.3,
// 		Normalisation: 3,
// 		Fold1:         0.5,
// 		Fold2:         1.5,
// 		AllIn:         0.95,
// 	}

// 	for g.IsGameOn() {

// 		g.NewRound(false)
// 		ShowCards(g.Rounds[len(g.Rounds)-1].Players)
// 		g.BetBlinds(false)

// 		numberOfPlayers := g.GetNumberOfPlayersAlive()
// 		turnToPlay := (g.SmallBlindTurn + 2) % numberOfPlayers
// 		for !g.Rounds[len(g.Rounds)-1].IsBetTurnOver(turnToPlay) {
// 			Play(&(g.Rounds[len(g.Rounds)-1]), turnToPlay, g.BigBlindValue, float64(g.TotalStack), coeff, false)
// 			turnToPlay = (turnToPlay + 1) % numberOfPlayers
// 		}
// 		for g.Rounds[len(g.Rounds)-1].PlayersAlive > 1 && len(g.Rounds[len(g.Rounds)-1].SharedCards) < 5 {
// 			switch len(g.Rounds[len(g.Rounds)-1].SharedCards) {
// 			case 0:
// 				fmt.Print("\n\tHere comes the flop...\n\n")
// 				g.Rounds[len(g.Rounds)-1].SharedCards = cards.Deal(&g.Rounds[len(g.Rounds)-1].Deck, 3)
// 				time.Sleep(time.Second * 1)
// 				fmt.Println(g.Rounds[len(g.Rounds)-1].SharedCards[0].ToString())
// 				time.Sleep(time.Second * 1)
// 				fmt.Println(g.Rounds[len(g.Rounds)-1].SharedCards[1].ToString())
// 				time.Sleep(time.Second * 1)
// 				fmt.Println(g.Rounds[len(g.Rounds)-1].SharedCards[2].ToString())
// 				time.Sleep(time.Second * 1)
// 			case 3:
// 				fmt.Print("\n\tHere comes the turn...\n\n")
// 				g.Rounds[len(g.Rounds)-1].SharedCards = append(g.Rounds[len(g.Rounds)-1].SharedCards, cards.Deal(&g.Rounds[len(g.Rounds)-1].Deck, 1)...)
// 				time.Sleep(time.Second * 1)
// 				fmt.Println(g.Rounds[len(g.Rounds)-1].SharedCards[3].ToString())
// 				time.Sleep(time.Second * 1)
// 				fmt.Print("\n\tCARDS:\n")
// 				cards.PrintDeck(g.Rounds[len(g.Rounds)-1].SharedCards)
// 			case 4:
// 				fmt.Print("\n\tHere comes the river...\n\n")
// 				g.Rounds[len(g.Rounds)-1].SharedCards = append(g.Rounds[len(g.Rounds)-1].SharedCards, cards.Deal(&g.Rounds[len(g.Rounds)-1].Deck, 1)...)
// 				time.Sleep(time.Second * 1)
// 				fmt.Println(g.Rounds[len(g.Rounds)-1].SharedCards[4].ToString())
// 				time.Sleep(time.Second * 1)
// 				fmt.Print("\n\tCARDS:\n")
// 				cards.PrintDeck(g.Rounds[len(g.Rounds)-1].SharedCards)
// 			}

// 			for i := range g.Rounds[len(g.Rounds)-1].Players {
// 				g.Rounds[len(g.Rounds)-1].Players[i].HasSpoken = false
// 			}
// 			turnToPlay := g.SmallBlindTurn
// 			for !g.Rounds[len(g.Rounds)-1].IsBetTurnOver(turnToPlay) {
// 				Play(&(g.Rounds[len(g.Rounds)-1]), turnToPlay, g.BigBlindValue, float64(g.TotalStack), coeff, false)
// 				turnToPlay = (turnToPlay + 1) % numberOfPlayers
// 			}
// 		}
// 		g.EndRound(false)
// 		time.Sleep(time.Second * 3)

// 	}
// 	fmt.Println("\tCongratulations", g.GetWinner(), "\n\n\tyou won!")
// }

// // func main() {
// //	initialProbabilty()
// //}

// func BotMode() {
// 	bot.InitialProbabilty()
// 	andrePoints := 0
// 	gomezPoints := 0

// 	coeffAndre := game.Coefficients{
// 		Check:         0.4,
// 		Raise:         0.2,
// 		Call:          0.1,
// 		Normalisation: 1,
// 		Fold1:         0.3,
// 		Fold2:         1.0,
// 		AllIn:         0.85,
// 	}

// 	coeffGomez := game.Coefficients{
// 		Check:         0.6, //  between 0.3 and 0.7
// 		Raise:         1.0, // between 0 and 1
// 		Call:          0.7, // between 0 and 1
// 		Normalisation: 1,   // not used
// 		Fold1:         0.8, // between 0.2 and 0.8
// 		Fold2:         5.0, // between 0.1 and 10
// 		AllIn:         1,   // between 0.7 and 1
// 	}
// 	first := true
// 	turn := 0
// 	difference := 0
// 	lastDifference := 0
// 	initialValue := float64(0.0)
// 	var updateDifference bool
// 	totalRounds := 0
// 	for math.Abs(coeffGomez.Fold1-coeffAndre.Fold1) > 0.1 || math.Abs(coeffGomez.Fold2-coeffAndre.Fold2) > 0.1 || math.Abs(coeffGomez.Call-coeffAndre.Call) > 0.1 || math.Abs(coeffGomez.Raise-coeffAndre.Raise) > 0.1 || math.Abs(coeffGomez.AllIn-coeffAndre.AllIn) > 0.1 || math.Abs(coeffGomez.Check-coeffAndre.Check) > 0.1 {

// 		difference = andrePoints - gomezPoints

// 		updateDifference = true
// 		if !first {
// 			if difference > 0 {
// 				if ((difference > lastDifference-20 && lastDifference > 0) || initialValue == 0) || totalRounds/(andrePoints+gomezPoints) < 20 {
// 					if turn == 0 {
// 						if initialValue == 0 {
// 							initialValue = coeffGomez.Raise
// 							coeffGomez.Raise = (initialValue + coeffAndre.Raise) / 2.0
// 						} else {
// 							if initialValue > coeffGomez.Raise {
// 								updateDifference = false
// 								coeffGomez.Raise = 2*initialValue - coeffGomez.Raise
// 							} else {
// 								turn++
// 								updateDifference = false
// 								coeffGomez.Raise = initialValue
// 								initialValue = 0
// 							}
// 						}
// 					}
// 					if turn == 1 {
// 						if initialValue == 0 {
// 							initialValue = coeffGomez.Call
// 							coeffGomez.Call = (initialValue + coeffAndre.Call) / 2.0
// 						} else {
// 							if initialValue > coeffGomez.Call {
// 								updateDifference = false
// 								coeffGomez.Call = 2*initialValue - coeffGomez.Call
// 							} else {
// 								turn++
// 								updateDifference = false
// 								coeffGomez.Call = initialValue
// 								initialValue = 0
// 							}
// 						}
// 					}
// 					if turn == 2 {
// 						if initialValue == 0 {
// 							initialValue = coeffGomez.Fold1
// 							coeffGomez.Fold1 = (initialValue + coeffAndre.Fold1) / 2.0
// 						} else {
// 							if initialValue > coeffGomez.Fold1 {
// 								updateDifference = false
// 								coeffGomez.Fold1 = 2*initialValue - coeffGomez.Fold1
// 							} else {
// 								turn++
// 								updateDifference = false
// 								coeffGomez.Fold1 = initialValue
// 								initialValue = 0
// 							}
// 						}
// 					}
// 					if turn == 3 {
// 						if initialValue == 0 {
// 							initialValue = coeffGomez.Fold2
// 							coeffGomez.Fold2 = (initialValue + coeffAndre.Fold2) / 2.0
// 						} else {
// 							if initialValue > coeffGomez.Fold2 {
// 								updateDifference = false
// 								coeffGomez.Fold2 = 2*initialValue - coeffGomez.Fold2
// 							} else {
// 								turn++
// 								updateDifference = false
// 								coeffGomez.Fold2 = initialValue
// 								initialValue = 0
// 							}
// 						}
// 					}
// 					if turn == 4 {
// 						if initialValue == 0 {
// 							initialValue = coeffGomez.AllIn
// 							coeffGomez.AllIn = (initialValue + coeffAndre.AllIn) / 2.0
// 						} else {
// 							if initialValue > coeffGomez.AllIn {
// 								updateDifference = false
// 								coeffGomez.AllIn = 2*initialValue - coeffGomez.AllIn
// 							} else {
// 								turn++
// 								updateDifference = false
// 								coeffGomez.AllIn = initialValue
// 								initialValue = 0
// 							}
// 						}
// 					}
// 					if turn == 5 {
// 						if initialValue == 0 {
// 							initialValue = coeffGomez.Check
// 							coeffGomez.Check = (initialValue + coeffAndre.Check) / 2.0
// 						} else {
// 							if initialValue > coeffGomez.Check {
// 								updateDifference = false
// 								coeffGomez.Check = 2*initialValue - coeffGomez.Check
// 							} else {
// 								turn = 0
// 								updateDifference = false
// 								coeffGomez.Check = initialValue
// 								initialValue = 0
// 							}
// 						}
// 					}
// 				} else {
// 					initialValue = 0
// 					if turn == 0 {
// 						fmt.Println("RAISE STABLE")
// 						turn++
// 					} else {
// 						if turn == 1 {
// 							fmt.Println("CALL STABLE")
// 							turn++
// 						} else {
// 							if turn == 2 {
// 								fmt.Println("FOLD1 STABLE")
// 								turn++
// 							} else {
// 								if turn == 3 {
// 									fmt.Println("FOLD2 STABLE")
// 									turn++
// 								} else {
// 									if turn == 4 {
// 										fmt.Println("ALLIN STABLE")
// 										turn++
// 									} else {
// 										fmt.Println("CHECK STABLE")
// 										turn = 0
// 									}
// 								}
// 							}
// 						}
// 					}
// 				}
// 			} else {
// 				if difference < 0 {
// 					if ((difference < lastDifference-8 && lastDifference < 0) || initialValue == 0) || totalRounds/(andrePoints+gomezPoints) < 20 {
// 						if turn == 0 {
// 							if initialValue == 0 {
// 								initialValue = coeffAndre.Raise
// 								coeffAndre.Raise = (initialValue + coeffGomez.Raise) / 2.0
// 							} else {
// 								if initialValue < coeffAndre.Raise {
// 									coeffAndre.Raise = 2*initialValue - coeffAndre.Raise
// 									updateDifference = false
// 								} else {
// 									turn++
// 									coeffAndre.Raise = initialValue
// 									updateDifference = false
// 									initialValue = 0
// 								}
// 							}
// 						}
// 						if turn == 1 {
// 							if initialValue == 0 {
// 								initialValue = coeffAndre.Call
// 								coeffAndre.Call = (initialValue + coeffGomez.Call) / 2.0
// 							} else {
// 								if initialValue < coeffAndre.Call {
// 									coeffAndre.Call = 2*initialValue - coeffAndre.Call
// 									updateDifference = false
// 								} else {
// 									turn++
// 									coeffAndre.Call = initialValue
// 									updateDifference = false
// 									initialValue = 0
// 								}
// 							}
// 						}
// 						if turn == 2 {
// 							if initialValue == 0 {
// 								initialValue = coeffAndre.Fold1
// 								coeffAndre.Fold1 = (initialValue + coeffGomez.Fold1) / 2.0
// 							} else {
// 								if initialValue < coeffAndre.Fold1 {
// 									coeffAndre.Fold1 = 2*initialValue - coeffAndre.Fold1
// 									updateDifference = false
// 								} else {
// 									turn++
// 									coeffAndre.Fold1 = initialValue
// 									updateDifference = false
// 									initialValue = 0
// 								}
// 							}
// 						}
// 						if turn == 3 {
// 							if initialValue == 0 {
// 								initialValue = coeffAndre.Fold2
// 								coeffAndre.Fold2 = (initialValue + coeffGomez.Fold2) / 2.0
// 							} else {
// 								if initialValue < coeffAndre.Fold2 {
// 									coeffAndre.Fold2 = 2*initialValue - coeffAndre.Fold2
// 									updateDifference = false
// 								} else {
// 									turn++
// 									coeffAndre.Fold2 = initialValue
// 									initialValue = 0
// 									updateDifference = false
// 								}
// 							}
// 						}
// 						if turn == 4 {
// 							if initialValue == 0 {
// 								initialValue = coeffAndre.AllIn
// 								coeffAndre.AllIn = (initialValue + coeffGomez.AllIn) / 2.0
// 							} else {
// 								if initialValue < coeffAndre.AllIn {
// 									coeffAndre.AllIn = 2*initialValue - coeffAndre.AllIn
// 									updateDifference = false
// 								} else {
// 									turn++
// 									coeffAndre.AllIn = initialValue
// 									initialValue = 0
// 									updateDifference = false
// 								}
// 							}
// 						}
// 						if turn == 5 {
// 							if initialValue == 0 {
// 								initialValue = coeffAndre.Check
// 								coeffAndre.Check = (initialValue + coeffGomez.Check) / 2.0
// 							} else {
// 								if initialValue < coeffAndre.Check {
// 									coeffAndre.Check = 2*initialValue - coeffAndre.Check
// 									updateDifference = false
// 								} else {
// 									turn = 0
// 									coeffAndre.Check = initialValue
// 									initialValue = 0
// 									updateDifference = false
// 								}
// 							}
// 						}
// 					} else {
// 						initialValue = 0
// 						if turn == 0 {
// 							fmt.Println("RAISE STABLE")
// 							turn++
// 						} else {
// 							if turn == 1 {
// 								fmt.Println("CALL STABLE")
// 								turn++
// 							} else {
// 								if turn == 2 {
// 									fmt.Println("FOLD1 STABLE")
// 									turn++
// 								} else {
// 									if turn == 3 {
// 										fmt.Println("FOLD2 STABLE")
// 										turn++
// 									} else {
// 										if turn == 4 {
// 											fmt.Println("ALLIN STABLE")
// 											turn++
// 										} else {
// 											fmt.Println("CHECK STABLE")
// 											turn = 0
// 										}
// 									}
// 								}
// 							}
// 						}
// 					}
// 				}

// 			}

// 			/*if andrePoints > gomezPoints {
// 				coeffGomez.raise = (coeffAndre.raise + coeffGomez.raise) / 2.0
// 				coeffGomez.call = (coeffAndre.call + coeffGomez.call) / 2.0
// 				coeffGomez.fold1 = (coeffAndre.fold1 + coeffGomez.fold1) / 2.0
// 				coeffGomez.fold2 = (coeffAndre.fold2 + coeffGomez.fold2) / 2.0
// 			} else {
// 				coeffAndre.raise = (coeffAndre.raise + coeffGomez.raise) / 2.0
// 				coeffAndre.call = (coeffAndre.call + coeffGomez.call) / 2.0
// 				coeffAndre.fold1 = (coeffAndre.fold1 + coeffGomez.fold1) / 2.0
// 				coeffAndre.fold2 = (coeffAndre.fold2 + coeffGomez.fold2) / 2.0
// 			}*/
// 		}
// 		andrePoints = 0
// 		gomezPoints = 0
// 		totalRounds = 0
// 		if updateDifference {
// 			lastDifference = difference
// 		}
// 		first = false

// 		fmt.Println("ANDRE VS GOMEZ, first in 300")
// 		fmt.Printf("%+v\n", coeffAndre)
// 		fmt.Printf("%+v\n", coeffGomez)
// 		for andrePoints < 300 && gomezPoints < 300 {
// 			//	fmt.Println("ANDRE", andrePoints, "-", gomezPoints, "GOMEZ")

// 			players := []game.Person{
// 				{
// 					Name:  "Andre",
// 					Stack: 1000,
// 					IsBot: true,
// 				},
// 				{
// 					Name:  "Gomez",
// 					Stack: 1000,
// 					IsBot: true,
// 				},
// 			}

// 			isBot := true
// 			g := game.NewGame(players)

// 			for g.IsGameOn() {

// 				g.NewRound(isBot)
// 				if !isBot {
// 					ShowCards(g.Rounds[len(g.Rounds)-1].Players)
// 				}
// 				g.BetBlinds(isBot)

// 				numberOfPlayers := g.GetNumberOfPlayersAlive()
// 				turnToPlay := (g.SmallBlindTurn + 2) % numberOfPlayers
// 				for !g.Rounds[len(g.Rounds)-1].IsBetTurnOver(turnToPlay) {
// 					if turnToPlay == 0 {
// 						Play(&(g.Rounds[len(g.Rounds)-1]), turnToPlay, g.BigBlindValue, float64(g.TotalStack), coeffAndre, isBot)
// 					} else {
// 						Play(&(g.Rounds[len(g.Rounds)-1]), turnToPlay, g.BigBlindValue, float64(g.TotalStack), coeffGomez, isBot)
// 					}

// 					turnToPlay = (turnToPlay + 1) % numberOfPlayers
// 				}
// 				for g.Rounds[len(g.Rounds)-1].PlayersAlive > 1 && len(g.Rounds[len(g.Rounds)-1].SharedCards) < 5 {
// 					switch len(g.Rounds[len(g.Rounds)-1].SharedCards) {
// 					case 0:
// 						g.Rounds[len(g.Rounds)-1].SharedCards = cards.Deal(&g.Rounds[len(g.Rounds)-1].Deck, 3)
// 						//printDeck(g.Rounds[len(g.Rounds)-1].SharedCards)
// 					case 3:
// 						g.Rounds[len(g.Rounds)-1].SharedCards = append(g.Rounds[len(g.Rounds)-1].SharedCards, cards.Deal(&g.Rounds[len(g.Rounds)-1].Deck, 1)...)
// 					//	printDeck(g.Rounds[len(g.Rounds)-1].SharedCards)
// 					case 4:
// 						g.Rounds[len(g.Rounds)-1].SharedCards = append(g.Rounds[len(g.Rounds)-1].SharedCards, cards.Deal(&g.Rounds[len(g.Rounds)-1].Deck, 1)...)
// 						//	printDeck(g.rounds[len(g.rounds)-1].SharedCards)
// 					}

// 					for i := range g.Rounds[len(g.Rounds)-1].Players {
// 						g.Rounds[len(g.Rounds)-1].Players[i].HasSpoken = false
// 					}
// 					turnToPlay := g.SmallBlindTurn
// 					for !g.Rounds[len(g.Rounds)-1].IsBetTurnOver(turnToPlay) {
// 						if turnToPlay == 0 {
// 							Play(&(g.Rounds[len(g.Rounds)-1]), turnToPlay, g.BigBlindValue, float64(g.TotalStack), coeffAndre, isBot)
// 						} else {
// 							Play(&(g.Rounds[len(g.Rounds)-1]), turnToPlay, g.BigBlindValue, float64(g.TotalStack), coeffGomez, isBot)
// 						}
// 						turnToPlay = (turnToPlay + 1) % numberOfPlayers
// 					}
// 				}
// 				g.EndRound(isBot)

// 			}
// 			if winner := g.GetWinnerIndex(); winner == 0 {
// 				andrePoints++
// 			} else {
// 				gomezPoints++
// 			}
// 			totalRounds += len(g.Rounds)
// 		}

// 		fmt.Println("ANDRE", andrePoints, "-", gomezPoints, "GOMEZ")
// 		fmt.Println("average rounds:", (totalRounds)/(andrePoints+gomezPoints))

// 	}

// }

// func ShowCards(sp []game.Player) {
// 	realPlayers := 0
// 	for _, p := range sp {
// 		if !p.IsBot {
// 			realPlayers++
// 		}
// 	}
// 	for _, p := range sp {
// 		if realPlayers > 1 {
// 			if p.IsBot {
// 				fmt.Print("\t", p.Name, ", your cards are:\n\n")
// 				fmt.Println(p.Cards[0].ToString(), " ", p.Cards[1].ToString())
// 			} else {
// 				c := ""
// 				for strings.ToLower(c) != "ok" {
// 					fmt.Print("\t", p.Name, " press OK to see your cards ")
// 					c = strings.ToLower(ReadFromTerminal())
// 				}
// 				fmt.Print("\t", p.Name, ", your cards are:\n\n")
// 				fmt.Println(p.Cards[0].ToString(), " ", p.Cards[1].ToString())
// 				c = ""
// 				for strings.ToLower(c) != "ok" {
// 					fmt.Print("\n\n\t", p.Name, " press OK to hide your cards ")
// 					c = strings.ToLower(ReadFromTerminal())
// 				}

// 			}
// 		} else {
// 			if realPlayers == 1 {
// 				if !p.IsBot {
// 					fmt.Print("\t", p.Name, ", your cards are:\n\n")
// 					fmt.Println(p.Cards[0].ToString(), " ", p.Cards[1].ToString())
// 				}
// 			} else {
// 				fmt.Print("\t", p.Name, ",'s cards are:\n\n")
// 				fmt.Println(p.Cards[0].ToString(), " ", p.Cards[1].ToString())
// 			}

// 		}
// 	}
// }

// func Play(r *game.Round, playerIndex int, bigBlind int, totalStacks float64, coeff game.Coefficients, botMode bool) {
// 	pp := &((*r).Players[playerIndex])
// 	p := (*pp)
// 	if p.IsAllIn && !botMode {
// 		fmt.Print("\t", p.Name, " is all in\n")
// 	}
// 	if !p.IsAllIn && !p.HasFolded {
// 		if !p.IsBot {
// 			fmt.Print("\t-> ", p.Name, " What do you want to do ? You have ", p.InitialStack-p.RoundBet, "\n")
// 			if (*r).MaxBet == p.RoundBet {
// 				fmt.Println("\t\tCheck (ch)")
// 			} else {
// 				fmt.Println("\t\tCall", (*r).MaxBet-p.RoundBet, "(ca)")
// 			}
// 			if p.InitialStack > (*r).MaxBet {
// 				fmt.Println("\t\tRaise (r)")
// 			}
// 			fmt.Print("\t\tFold (f) ")
// 			validInput := false
// 			for !validInput {
// 				validInput = true
// 				c := strings.ToLower(ReadFromTerminal())
// 				switch c {
// 				case "check", "ch":
// 					if (*r).MaxBet != p.RoundBet {
// 						fmt.Print("\t", p.Name, ", checks\n\n")
// 						validInput = false
// 					}
// 				case "call", "ca":
// 					if (*r).MaxBet != p.RoundBet {
// 						fmt.Print("\t", p.Name, " calls\n\n")
// 						(*r).Bet(playerIndex, (*r).MaxBet-p.RoundBet, false)
// 					}
// 				case "fold", "f":
// 					(*pp).HasFolded = true
// 					(*r).PlayersAlive--
// 					fmt.Print("\t", p.Name, " folds\n\n")
// 				case "raise", "r":
// 					if p.InitialStack > (*r).MaxBet {
// 						valueRaise := 0
// 						for valueRaise == 0 {
// 							if p.InitialStack-(*r).MaxBet > bigBlind {
// 								fmt.Print("\t---> How much do you want to raise? (", bigBlind, "-", p.InitialStack-(*r).MaxBet, ")")
// 								text := ReadFromTerminal()
// 								number, err := strconv.Atoi(text)
// 								if err == nil {
// 									if number >= bigBlind {
// 										valueRaise = number
// 										(*r).Bet(playerIndex, (*r).MaxBet-p.RoundBet+number, false)
// 									}
// 								}
// 							} else {
// 								(*r).Bet(playerIndex, p.InitialStack-p.RoundBet, false)
// 							}
// 						}
// 						fmt.Print("\t", p.Name, " raises ", valueRaise, "\n\n")
// 					}
// 				default:
// 					validInput = false
// 				}
// 				if !validInput {
// 					fmt.Println("Your choice has not been understood")
// 				}
// 			}
// 		} else {
// 			var decision float64

// 			if (*pp).HasSpoken {
// 				decision = (*pp).Decision
// 			} else {
// 				odds := odds.CalculateOdds(p.Cards, r.SharedCards)
// 				//fmt.Println("odds:", odds)
// 				/*	var exp float64
// 					if len(r.sharedCards) == 0 {
// 						exp = expoFunction(odds, 3.0)
// 					} else {
// 						exp = expoFunction(odds, 2.0)
// 					}
// 					//		fmt.Println("expo:", exp)
// 					nlr := normaleLawRandomization(1, 0)*/
// 				//		fmt.Println("normalization:", nlr)
// 				//decision = exp * nlr
// 				//fmt.Println("decision before random:", odds)
// 				decision = odds * NormaleLawRandomization(0.7, 0)
// 				(*pp).Decision = decision
// 			}

// 			//fmt.Println("decision after random:", decision)
// 			if (*r).MaxBet == p.RoundBet {
// 				if decision < coeff.Check {
// 					if !botMode {
// 						fmt.Print("\t", p.Name, " checks\n\n")
// 					}
// 				} else {
// 					raiseValue := cards.MaxInt(int(math.Ceil(ExpoFunction(decision, 3.0)*(totalStacks/float64(len(r.Players)))*coeff.Raise*NormaleLawRandomization(1, 0))), bigBlind)
// 					raiseValue -= raiseValue % (bigBlind / 2)
// 					(*r).Bet(playerIndex, raiseValue, botMode)
// 					if !botMode {
// 						fmt.Print("\t", p.Name, " raises ", raiseValue, "\n\n")
// 					}

// 				}
// 			} else {
// 				toCallRatio := float64((*r).MaxBet-p.RoundBet) / (totalStacks * math.Exp(float64(p.RoundBet)/(totalStacks/float64(len(r.Players)))))

// 				limit1 := coeff.Fold1
// 				limit2 := RaiseExpo(toCallRatio * coeff.Fold2)
// 				//fmt.Println("tocallratio:", limit2)
// 				//fmt.Println("decision:", decision)
// 				if decision < limit1 || (decision < limit2) /* && (float64(p.initialStack-p.roundBet) > 0.05*totalStacks*/ {
// 					(*pp).HasFolded = true
// 					(*r).PlayersAlive--
// 					if !botMode {
// 						fmt.Print("\t", p.Name, " folds\n\n")
// 					}

// 				} else {
// 					overFold := decision - limit2
// 					if overFold < coeff.Call {
// 						(*r).Bet(playerIndex, (*r).MaxBet-p.RoundBet, botMode)
// 						if !botMode {
// 							fmt.Print("\t", p.Name, " calls\n\n")
// 						}

// 					} else {
// 						var raiseValue int
// 						if decision > coeff.AllIn {
// 							raiseValue = p.InitialStack - p.RoundBet
// 						} else {
// 							raiseValue = cards.MaxInt(int(math.Ceil(ExpoFunction(overFold, 3.0)*(totalStacks/float64(len(r.Players)))*coeff.Raise*NormaleLawRandomization(0.5, 0))), bigBlind)
// 							raiseValue -= raiseValue % (bigBlind / 2)
// 						}

// 						(*r).Bet(playerIndex, (*r).MaxBet-p.RoundBet+raiseValue, botMode)
// 						if !botMode {
// 							fmt.Print("\t", p.Name, " raises ", raiseValue, "\n\n")
// 						}

// 					}
// 				}
// 			}
// 		}
// 		(*pp).HasSpoken = true

// 	}

// }

// func ExpoFunction(odds float64, index float64) float64 {
// 	return math.Exp(index * (odds - 1.0))
// }

// func RaiseExpo(x float64) float64 {
// 	return (1.5 - 0.5*x) * math.Exp(x-1.0)
// }

// func NormaleLawRandomization(risk float64, offset float64) float64 {
// 	source := rand.NewSource(time.Now().UnixNano())
// 	r := rand.New(source)
// 	riskInt := int(math.Ceil(risk * 100))
// 	randNumber := r.Intn(riskInt*2) - riskInt
// 	x := float64(randNumber)/100 + offset
// 	if x < 0 {
// 		return math.Exp(-0.5 * x * x)
// 	}

// 	return 2.0 - math.Exp(-0.5*x*x)

// }

// func RatioStacks(stack float64, opponentStack float64) float64 {
// 	return (math.Log10(stack/opponentStack) + 2.0) / 2.0
// }

// func RatioToWintoBet(amountToBet float64, pot float64) float64 {
// 	return pot / amountToBet
// }

// func ReadFromTerminal() string {
// 	reader := bufio.NewReader(os.Stdin)
// 	text, _ := reader.ReadString('\n')
// 	text = strings.Replace(text, "\n", "", -1)
// 	text = strings.Replace(text, "\r", "", -1)

// 	return text
// }
