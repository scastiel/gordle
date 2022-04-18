package gordle

import (
	"strings"

	"fyne.io/fyne/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type AppState struct {
	game         Game
	guesses      []Word
	feedbacks    []Feedback
	currentWord  string
	errorMessage string
	letterColors map[string]Color
	aboutWindow  *fyne.Window
}

func NewAppState() *AppState {
	return &AppState{
		game:         NewGame(),
		guesses:      make([]Word, 0),
		feedbacks:    make([]Feedback, 0),
		currentWord:  "",
		errorMessage: "",
		letterColors: make(map[string]Color),
		aboutWindow:  nil,
	}
}

func (state AppState) typeLetter(letter string) AppState {
	if state.game.State == Started && len(state.currentWord) < 5 {
		state.currentWord += letter
	}
	return state
}

func (state AppState) backspace() AppState {
	if state.game.State == Started && len(state.currentWord) > 0 {
		state.currentWord = state.currentWord[0 : len(state.currentWord)-1]
	}
	return state
}

func (state AppState) enter() AppState {
	if state.game.State == Started {
		guess, err := NewWord(state.currentWord)
		if err != nil {
			return state.setError(err)
		}
		var feedback Feedback
		state.game, feedback, err = state.game.Guess(guess)
		if err != nil {
			return state.setError(err)
		}
		state.guesses = append(state.guesses, guess)
		state.feedbacks = append(state.feedbacks, feedback)

		for i, letter := range strings.Split(guess.Letters, "") {
			color := feedback.Colors[i]
			switch color {
			case Green:
				state.letterColors[letter] = Green
			case Yellow:
				if state.letterColors[letter] != Green {
					state.letterColors[letter] = Yellow
				}
			case Grey:
				if state.letterColors[letter] != Green && state.letterColors[letter] != Yellow {
					state.letterColors[letter] = Grey
				}
			}
		}

		state.currentWord = ""
	} else {
		state = *NewAppState()
	}
	return state
}

func (state AppState) resetError() AppState {
	state.errorMessage = ""
	return state
}

func (state AppState) setError(err error) AppState {
	state.errorMessage = cases.Title(language.English).String(err.Error())
	return state
}

func (state AppState) setAboutWindow(window *fyne.Window) AppState {
	state.aboutWindow = window
	return state
}
