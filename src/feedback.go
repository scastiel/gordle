package gordle

import (
	"errors"
)

type Color string

const (
	Grey   Color = "grey"
	Yellow Color = "yellow"
	Green  Color = "green"
)

type Feedback struct {
	Colors []Color
}

func (feedback *Feedback) IsWin() bool {
	return feedback.Colors[0] == Green &&
		feedback.Colors[1] == Green &&
		feedback.Colors[2] == Green &&
		feedback.Colors[3] == Green &&
		feedback.Colors[4] == Green
}

func NewFeedback(colors []Color) (Feedback, error) {
	if len(colors) != 5 {
		return Feedback{}, errors.New("need 5 colors")
	}
	for _, color := range colors {
		if color != Grey && color != Yellow && color != Green {
			return Feedback{}, errors.New("invalid color")
		}
	}
	return Feedback{colors}, nil
}

func FromGuess(guess Word, solution Word) Feedback {
	colors := []Color{Grey, Grey, Grey, Grey, Grey}
	excluded := []bool{false, false, false, false, false}

	for i := range colors {
		if guess.Letters[i] == solution.Letters[i] {
			colors[i] = Green
			excluded[i] = true
		}
	}

	for i := range colors {
		if colors[i] != Grey {
			continue
		}
		for j := range colors {
			if excluded[j] {
				continue
			}
			if guess.Letters[i] == solution.Letters[j] {
				colors[i] = Yellow
				excluded[j] = true
				break
			}
		}
	}

	feedback, _ := NewFeedback(colors)
	return feedback
}
