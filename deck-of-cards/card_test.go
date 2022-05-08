package deck

import (
	"fmt"
	"testing"
)

func ExampleCard() {
	fmt.Println(Card{Rank: Ace, Suit: Spade})
	fmt.Println(Card{Rank: Ace, Suit: Heart})
	fmt.Println(Card{Suit: Joker})

	// Output:
	// Ace of Spades
	// Ace of Hearts
	// Joker
}

func TestNew(t *testing.T) {
	cards := New()
	if len(cards) != 13*4 {
		t.Error("Wrong number of cards in a new deck.")
	}
}

func TestDefaultSort(t *testing.T) {
	cards := New(DefaultSort)
	card := Card{Suit: Spade, Rank: Ace}
	if cards[0] != card {
		t.Errorf("expected: %s, actual: %s", card, cards[0])
	}
}

func TestSort(t *testing.T) {
	cards := New(Sort(Less))
	card := Card{Suit: Spade, Rank: Ace}
	if cards[0] != card {
		t.Errorf("expected: %s, actual: %s", card, cards[0])
	}
}

func TestShuffle(t *testing.T) {
	shuffledDeck := New(Shuffle(52))
	var cards = [...]Card{
		{Rank: Rank(Four), Suit: Suit(Club)},
		{Rank: Rank(Four), Suit: Suit(Diamond)},
	}
	for i, c := range cards {
		if shuffledDeck[i] != c {
			t.Errorf("expected %s, got %s", c, shuffledDeck[i])
		}
	}
}

func TestFilter(t *testing.T) {
	predicate := func(card Card) bool {
		return card.Rank == Two || card.Rank == Three
	}
	cards := New(Filter(predicate))
	for _, c := range cards {
		if c.Rank == Two || c.Rank == Three {
			t.Errorf("twos and threes shouldn't be present in the deck")
		}
	}
}

func TestJokers(t *testing.T) {
	cards := New(Jokers(4))
	count := 0
	for _, c := range cards {
		if c.Suit == Joker {
			count++
		}
	}
	if count != 4 {
		t.Errorf("expected %d, got: %d", 4, count)
	}
}

func TestDeck(t *testing.T) {
	decks := New(Deck(3))
	if len(decks) != 13*4*3 {
		t.Errorf("expected %d, got: %d", 13*4*3, len(decks))
	}
}
