package blackjack

import (
	"fmt"

	"deck"
)

type AI interface {
	Bet() int
	Play(hand []deck.Card, dealer deck.Card) Move
	Results(hand [][]deck.Card, dealer []deck.Card)
}

type dealerAI struct{}

func (ai dealerAI) Bet() int {
	// noop
	return 1
}

func (ai dealerAI) Play(hand []deck.Card, dealer deck.Card) Move {
	dScore := Score(hand...)
	if dScore <= 16 || (dScore == 17 && Soft(hand...)) {
		return MoveHit
	}
	return MoveStand
}

func (ai dealerAI) Results(hand [][]deck.Card, dealer []deck.Card) {}

func HumanAI() AI {
	return humanAI{}
}

type humanAI struct{}

func (ai humanAI) Bet() int {
	var input int
	fmt.Println("How much do you want to bet?")
	fmt.Scanf("%d\n", &input)
	return input
}

func (ai humanAI) Play(hand []deck.Card, dealer deck.Card) Move {
	for {
		fmt.Println("Player:", Hand(hand))
		fmt.Println("Dealer:", dealer)
		fmt.Println("What will you do? (h)it, (s)tand")
		var input string
		fmt.Scanf("%s\n", &input)
		switch input {
		case "h":
			return MoveHit
		case "s":
			return MoveStand
		default:
			fmt.Println("Invalid option:", input)
		}
	}
}

func (ai humanAI) Results(hand [][]deck.Card, dealer []deck.Card) {
	fmt.Println("==FINAL HANDS==")
	var player Hand
	for _, h := range hand {
		for _, v := range h {
			player = append(player, v)
		}
	}
	fmt.Println("Player: ", player, "\nScore: ", Score(player...))
	fmt.Println("Dealer:", Hand(dealer), "\nScore: ", Score(dealer...))
}
