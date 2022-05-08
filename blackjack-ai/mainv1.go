package main

import (
	"deck"
	"fmt"
	"strings"
)

type Hand []deck.Card

type State int8

const (
	StatePlayerTurn State = iota
	StateDealerTurn
	StateHandOver
)

type GameState struct {
	deck   []deck.Card
	state  State
	player Hand
	dealer Hand
}

func (gs *GameState) getCurrentPlayer() *Hand {
	switch gs.state {
	case StatePlayerTurn:
		return &gs.player
	case StateDealerTurn:
		return &gs.dealer
	default:
		panic("invalid state")
	}
}

func clone(gs GameState) GameState {
	ret := GameState{
		deck:   make([]deck.Card, len(gs.deck)),
		state:  gs.state,
		player: make(Hand, len(gs.player)),
		dealer: make(Hand, len(gs.dealer)),
	}
	copy(ret.deck, gs.deck)
	copy(ret.player, gs.player)
	copy(ret.dealer, gs.dealer)
	return ret
}

func Shuffle(gs GameState) GameState {
	ret := clone(gs)
	ret.deck = deck.New(deck.Deck(3), deck.Shuffle())
	return ret
}

func Deal(gs GameState) GameState {
	ret := clone(gs)
	ret.player = make(Hand, 0, 5)
	ret.dealer = make(Hand, 0, 5)
	var card deck.Card
	for i := 0; i < 2; i++ {
		card, ret.deck = draw(ret.deck)
		ret.player = append(ret.player, card)
		card, ret.deck = draw(ret.deck)
		ret.dealer = append(ret.dealer, card)
	}
	ret.state = StatePlayerTurn
	return ret
}

func Hit(gs GameState) GameState {
	ret := clone(gs)
	hand := ret.getCurrentPlayer()
	var card deck.Card
	card, ret.deck = draw(gs.deck)
	*hand = append(*hand, card)
	if hand.Score() > 21 {
		return Stand(ret)
	}
	return ret
}

func Stand(gs GameState) GameState {
	ret := clone(gs)
	ret.state += 1
	return ret
}

func EndHand(gs GameState) GameState {
	ret := clone(gs)
	pScore, dScore := ret.player.Score(), ret.dealer.Score()
	fmt.Println("==Final Hands==")
	fmt.Println("Player: ", gs.player, "\nScore: ", pScore)
	fmt.Println("Dealer: ", gs.dealer, "\nScore: ", dScore)

	switch {
	case pScore > 21:
		fmt.Print("You busted")
	case dScore > 21:
		fmt.Print("Dealer busted")
	case pScore > dScore:
		fmt.Print("You win")
	case dScore > pScore:
		fmt.Print("You lose")
	case dScore == pScore:
		fmt.Print("Draw")
	}
	fmt.Println()
	ret.player = nil
	ret.dealer = nil
	return ret
}

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

func (h Hand) String() string {
	ret := make([]string, len(h))
	for i := range h {
		ret[i] = h[i].String()
	}
	return strings.Join(ret, ", ")
}

func (h Hand) DealerString() string {
	return h[0].String() + ", **HIDDEN**"
}

func (h Hand) MinScore() int {
	score := 0
	for _, c := range h {
		score += min(int(c.Rank), 10)
	}
	return score
}

func (h Hand) Score() int {
	minScore := h.MinScore()
	if minScore > 11 {
		return minScore
	}

	for _, c := range h {
		if c.Rank == deck.Ace {
			return minScore + 10
		}
	}
	return minScore
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// func main() {
// 	var gs GameState
// 	gs = Shuffle(gs)
// 	for i := 0; i < 3; i++ {
// 		gs = Deal(gs)
// 		var input string
// 		for gs.state == StatePlayerTurn {
// 			fmt.Println("Player: ", gs.player)
// 			fmt.Println("Dealer: ", gs.dealer.DealerString())
// 			fmt.Println("What will you do? (h)it or (s)tand")
// 			fmt.Scanf("%s\n", &input)
// 			switch input {
// 			case "h":
// 				gs = Hit(gs)
// 			case "s":
// 				gs = Stand(gs)
// 			default:
// 				fmt.Println("Invalid input", input)
// 			}
// 		}
// 		for gs.state == StateDealerTurn {
// 			if gs.dealer.Score() <= 16 || (gs.dealer.MinScore() != 17 && gs.dealer.Score() == 17) {
// 				gs = Hit(gs)
// 			} else {
// 				gs = Stand(gs)
// 			}
// 		}
// 		gs = EndHand(gs)
// 	}
// }
