package blackjack

import (
	"fmt"
	"strings"

	"deck"
)

const (
	statePlayerTurn state = iota
	stateDealerTurn
	stateHandOver
)

type state int8

type Options struct {
	Decks           int
	Hands           int
	BlackJackPayout float64
}

func validateOptions(opts *Options) {
	if opts.Decks <= 0 {
		opts.Decks = 3
	}
	if opts.Hands <= 0 {
		opts.Hands = 2
	}
	if opts.BlackJackPayout <= 0 {
		opts.BlackJackPayout = 1.5
	}
}

func New(opts Options) Game {
	validateOptions(&opts)
	return Game{
		state:           statePlayerTurn,
		dealerAI:        dealerAI{},
		balance:         0,
		nDecks:          opts.Decks,
		nHands:          opts.Hands,
		blackJackPayout: opts.BlackJackPayout,
	}
}

type Game struct {
	// unexported fields
	deck            []deck.Card
	state           state
	player          []deck.Card
	dealer          []deck.Card
	dealerAI        AI
	balance         int
	nDecks          int
	nHands          int
	blackJackPayout float64
	playerBet       int
}

func (g *Game) currentHand() *[]deck.Card {
	switch g.state {
	case statePlayerTurn:
		return &g.player
	case stateDealerTurn:
		return &g.dealer
	default:
		panic("it isn't currently any player's turn")
	}
}

func bet(g *Game, ai AI) {
	g.playerBet = ai.Bet()
}

func deal(g *Game) {
	g.player = make([]deck.Card, 0, 5)
	g.dealer = make([]deck.Card, 0, 5)
	var card deck.Card
	for i := 0; i < 2; i++ {
		card, g.deck = draw(g.deck)
		g.player = append(g.player, card)
		card, g.deck = draw(g.deck)
		g.dealer = append(g.dealer, card)
	}
	g.state = statePlayerTurn
}

func (g *Game) Play(ai AI) int {
	g.deck = deck.New(deck.Deck(3), deck.Shuffle())
	for i := 0; i < 2; i++ {
		bet(g, ai)
		deal(g)
		for g.state == statePlayerTurn {
			hand := make([]deck.Card, len(g.player))
			copy(hand, g.player)
			move := ai.Play(hand, g.dealer[0])
			move(g)
		}

		for g.state == stateDealerTurn {
			hand := make([]deck.Card, len(g.dealer))
			copy(hand, g.dealer)
			move := g.dealerAI.Play(hand, g.dealer[0])
			move(g)
		}

		endRound(g, ai)
	}
	return g.balance
}

type Move func(*Game)

func MoveHit(g *Game) {
	hand := g.currentHand()
	var card deck.Card
	card, g.deck = draw(g.deck)
	*hand = append(*hand, card)
	if Score(*hand...) > 21 {
		MoveStand(g)
	}
}

func MoveStand(g *Game) {
	g.state++
}

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

func endRound(g *Game, ai AI) {
	pScore, dScore := Score(g.player...), Score(g.dealer...)
	winnings := g.playerBet
	switch {
	case pScore > 21:
		winnings = -1 * winnings
	case dScore > 21:
		// win
	case pScore > dScore:
		// win
	case dScore > pScore:
		winnings = -1 * winnings
	case dScore == pScore:
		winnings = 0
	}
	g.balance += winnings
	fmt.Println()
	ai.Results([][]deck.Card{g.player}, g.dealer)
	g.player = nil
	g.dealer = nil
}

// Score will take in a hand of cards and return the best blackjack score
// possible with the hand.
func Score(hand ...deck.Card) int {
	minScore := minScore(hand...)
	if minScore > 11 {
		return minScore
	}
	for _, c := range hand {
		if c.Rank == deck.Ace {
			// ace is currently worth 1, and we are changing it to be worth 11
			// 11 - 1 = 10
			return minScore + 10
		}
	}
	return minScore
}

// Soft returns true if the score of a hand is a soft score - that is if an ace
// is being counted as 11 points.
func Soft(hand ...deck.Card) bool {
	minScore := minScore(hand...)
	score := Score(hand...)
	return minScore != score
}

func minScore(hand ...deck.Card) int {
	score := 0
	for _, c := range hand {
		score += min(int(c.Rank), 10)
	}
	return score
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type Hand []deck.Card

func (h Hand) String() string {
	ret := make([]string, len(h))
	for i := range h {
		ret[i] = h[i].String()
	}
	return strings.Join(ret, ", ")
}
