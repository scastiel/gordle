package gordle

import (
	"errors"
	"math/rand"
	"regexp"
)

type Word struct {
	Letters string
}

func isValidWord(letters string) bool {
	for _, word := range validWords {
		if word == letters {
			return true
		}
	}
	return false
}

func NewWord(letters string) (Word, error) {
	regex := regexp.MustCompile("[A-Z]{5}")
	if !regex.MatchString(letters) {
		return Word{}, errors.New("invalid word")
	}
	if !isValidWord(letters) {
		return Word{}, errors.New("unknown word")
	}
	return Word{Letters: letters}, nil
}

func RandomWord() Word {
	letters := validWords[rand.Intn(len(validWords))]
	word, _ := NewWord(letters)
	return word
}
