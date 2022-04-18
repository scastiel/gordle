package gordle

import (
	"errors"
)

type State string

const (
	Started State = "started"
	Won     State = "won"
	Lost    State = "lost"
)

const maxTries int = 6

type Game struct {
	State    State
	Tries    int
	Solution Word
}

func NewGame() Game {
	return Game{State: "started", Tries: 0, Solution: RandomWord()}
}

func (game *Game) Guess(guess Word) (Game, Feedback, error) {
	if game.State != Started {
		return *game, Feedback{}, errors.New("Game already finished")
	}
	feedback := FromGuess(guess, game.Solution)
	if feedback.IsWin() {
		return Game{State: Won, Tries: game.Tries + 1, Solution: game.Solution}, feedback, nil
	}
	if game.Tries == maxTries-1 {
		return Game{State: Lost, Tries: game.Tries + 1, Solution: game.Solution}, feedback, nil
	}
	return Game{State: game.State, Tries: game.Tries + 1, Solution: game.Solution}, feedback, nil
}
