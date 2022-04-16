package gordle

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGameIsStarted(t *testing.T) {
	game := NewGame()
	assert.Equal(t, Started, game.State)
	assert.Equal(t, 0, game.Tries)
	assert.NotNil(t, game.Solution)
}

func TestGuessOnFinishedGameRaisesError(t *testing.T) {
	game := Game{Solution: RandomWord(), Tries: maxTries, State: Won}
	guess, _ := NewWord("HELLO")
	_, _, err := Guess(game, guess)
	assert.Error(t, errors.New("Game already finished"), err)
}

func TestGuessWithWrongAtLastTryLoses(t *testing.T) {
	game := Game{Solution: Word{Letters: "NAVAL"}, Tries: maxTries - 1, State: Started}
	game, feedback, err := Guess(game, Word{Letters: "HELLO"})
	assert.Equal(t, State(Lost), game.State)
	assert.Equal(t, feedback, Feedback{Colors: []Color{Grey, Grey, Yellow, Grey, Grey}})
	assert.Nil(t, err)
}

func TestGuessWithWrongNotAtLastTryContinues(t *testing.T) {
	game := Game{Solution: Word{Letters: "NAVAL"}, Tries: 1, State: Started}
	game, feedback, err := Guess(game, Word{Letters: "HELLO"})
	assert.Equal(t, 2, game.Tries)
	assert.Equal(t, State(Started), game.State)
	assert.Equal(t, feedback, Feedback{Colors: []Color{Grey, Grey, Yellow, Grey, Grey}})
	assert.Nil(t, err)
}

func TestGuessWithSolutionWins(t *testing.T) {
	game := Game{Solution: Word{Letters: "NAVAL"}, Tries: 1, State: Started}
	game, feedback, err := Guess(game, Word{Letters: "NAVAL"})
	assert.Equal(t, 2, game.Tries)
	assert.Equal(t, State(Won), game.State)
	assert.Equal(t, feedback, Feedback{Colors: []Color{Green, Green, Green, Green, Green}})
	assert.Nil(t, err)
}
