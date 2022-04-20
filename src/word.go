package gordle

import (
	"errors"
	"math/rand"
	"regexp"
	"sort"
)

type Word struct {
	Letters string
}

func init() {
	if !sort.StringsAreSorted(validWords) {
		sort.Strings(validWords)
	}
}

func isValidWord(letters string) bool {
	i := sort.SearchStrings(validWords, letters)
	return i < len(validWords) && validWords[i] == letters
}

var validWordRegex = regexp.MustCompile("[A-Z]{5}")

func NewWord(letters string) (Word, error) {
	if !validWordRegex.MatchString(letters) {
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
