// go:generate stringer -type=Rank,Suit
package deck

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type Suit uint8
type Rank uint8

type Card struct {
	Rank
	Suit
}

const (
	Spade Suit = iota
	Club
	Diamond
	Heart
	Joker
)

var suits = []Suit{Spade, Club, Diamond, Heart}

const (
	_ Rank = iota
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

const (
	minRank = Ace
	maxRank = King
)

func (c Card) String() string {
	if c.Suit == Joker {
		return fmt.Sprint(c.Suit.String())
	}
	return fmt.Sprintf("%s of %ss", c.Rank.String(), c.Suit.String())
}

func DefaultSort(cards []Card) []Card {
	sort.Slice(cards, Less(cards))
	return cards
}

func Sort(less func(cards []Card) func(i, j int) bool) func(cards []Card) []Card {
	return func(cards []Card) []Card {
		sort.Slice(cards, less(cards))
		return cards
	}
}

func Less(cards []Card) func(i, j int) bool {
	return func(i, j int) bool {
		return absRank(cards[i]) < absRank(cards[j])
	}
}

func absRank(card Card) int {
	return int(card.Suit)*int(maxRank) + int(card.Rank)
}

func New(opts ...func(c []Card) []Card) []Card {
	var cards []Card

	for _, suit := range suits {
		for rank := minRank; rank <= maxRank; rank++ {
			cards = append(cards, Card{Rank: rank, Suit: suit})
		}
	}

	for _, opt := range opts {
		cards = opt(cards)
	}

	return cards
}

func Shuffle(seed ...int64) func(cards []Card) []Card {
	if len(seed) == 0 {
		seed = append(seed, time.Now().Unix())
	}
	return func(cards []Card) []Card {
		var ret = make([]Card, len(cards))
		r := rand.New(rand.NewSource(seed[0]))
		perm := r.Perm(len(cards))
		for i, j := range perm {
			ret[i] = cards[j]
		}
		return ret
	}
}

func Jokers(n int) func([]Card) []Card {
	return func(cards []Card) []Card {
		for n > 0 {
			cards = append(cards, Card{
				Suit: Joker,
				Rank: Rank(n),
			})
			n -= 1
		}
		return cards
	}
}

func Filter(f func(card Card) bool) func(cards []Card) []Card {
	return func(cards []Card) []Card {
		var ret []Card
		for _, card := range cards {
			if !f(card) {
				ret = append(ret, card)
			}
		}
		return ret
	}
}

func Deck(n int) func(cards []Card) []Card {
	return func(cards []Card) []Card {
		var ret []Card
		for n > 0 {
			ret = append(ret, cards...)
			n -= 1
		}
		return ret
	}
}
