package gordle

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWordReturnsWordWhenValid(t *testing.T) {
	word, err := NewWord("CIGAR")
	assert.Equal(t, Word{Letters: "CIGAR"}, word)
	assert.Nil(t, err)
}

func TestWordReturnsErrorWhenInvalid(t *testing.T) {
	_, err := NewWord("invalid")
	assert.Error(t, errors.New("invalid word"), err)
}

func TestWordReturnsErrorWhenUnknown(t *testing.T) {
	_, err := NewWord("AAAAA")
	assert.Error(t, errors.New("unknown word"), err)
}

func TestRandomWordReturnsWord(t *testing.T) {
	w := RandomWord()
	assert.Equal(t, 5, len(w.Letters))
}
