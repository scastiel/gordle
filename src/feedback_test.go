package gordle

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFeedbackReturnsFeedbackWithValidColors(t *testing.T) {
	colors := []Color{Grey, Green, Yellow, Grey, Grey}
	f, err := NewFeedback(colors)
	assert.Equal(t, Feedback{colors}, f)
	assert.Nil(t, err)
}

func TestFeedbackReturnsErrorWithWrongNumberOfColors(t *testing.T) {
	colors := []Color{Grey, Green, Yellow, Grey}
	_, err := NewFeedback(colors)
	assert.Error(t, errors.New("need 5 colors"), err)
}

func TestFeedbackReturnsErrorWithInvalidColors(t *testing.T) {
	colors := []Color{Grey, Green, Yellow, Grey, "invalid"}
	_, err := NewFeedback(colors)
	assert.Error(t, errors.New("invalid color"), err)
}

func TestFromGuessReturnsAllGreenForSameWords(t *testing.T) {
	solution, _ := NewWord("NAVAL")
	guess, _ := NewWord("NAVAL")
	feedback := FromGuess(solution, guess)
	assert.Equal(t, Feedback{Colors: []Color{Green, Green, Green, Green, Green}}, feedback)
}

func TestFromGuessReturnsAllGreyForTotallyDifferentWords(t *testing.T) {
	solution, _ := NewWord("CIGAR")
	guess, _ := NewWord("HUMPH")
	feedback := FromGuess(solution, guess)
	assert.Equal(t, Feedback{Colors: []Color{Grey, Grey, Grey, Grey, Grey}}, feedback)
}

func TestFromGuessReturnsMixedColorsForSimilarWords(t *testing.T) {
	solution, _ := NewWord("HELLO")
	guess, _ := NewWord("FILTH")
	feedback := FromGuess(guess, solution)
	assert.Equal(t, Feedback{Colors: []Color{Grey, Grey, Green, Grey, Yellow}}, feedback)
}
