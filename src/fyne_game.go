package gordle

import (
	c "image/color"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
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
}

func (state *AppState) typeLetter(letter string) {
	if state.game.State == Started && len(state.currentWord) < 5 {
		state.currentWord += letter
	}
}

func (state *AppState) backspace() {
	if state.game.State == Started && len(state.currentWord) > 0 {
		state.currentWord = state.currentWord[0 : len(state.currentWord)-1]
	}
}

func (state *AppState) enter() error {
	if state.game.State == Started {
		guess, err := NewWord(state.currentWord)
		if err != nil {
			return err
		}
		var feedback Feedback
		state.game, feedback, err = Guess(state.game, guess)
		if err != nil {
			return err
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
	}
	return nil
}

func StartFyneGame() {
	state := AppState{
		game:         NewGame(),
		guesses:      make([]Word, 0),
		feedbacks:    make([]Feedback, 0),
		currentWord:  "",
		errorMessage: "",
		letterColors: make(map[string]Color),
	}

	app := app.New()
	icon, _ := fyne.LoadResourceFromPath("assets/AppIcon.svg")
	app.SetIcon(icon)
	window := app.NewWindow("Gordle")
	window.SetFixedSize(true)
	render(&state, window)

	mappings := map[fyne.KeyName]string{
		fyne.KeyA: "A",
		fyne.KeyB: "B",
		fyne.KeyC: "C",
		fyne.KeyD: "D",
		fyne.KeyE: "E",
		fyne.KeyF: "F",
		fyne.KeyG: "G",
		fyne.KeyH: "H",
		fyne.KeyI: "I",
		fyne.KeyJ: "J",
		fyne.KeyK: "K",
		fyne.KeyL: "L",
		fyne.KeyM: "M",
		fyne.KeyN: "N",
		fyne.KeyO: "O",
		fyne.KeyP: "P",
		fyne.KeyQ: "Q",
		fyne.KeyR: "R",
		fyne.KeyS: "S",
		fyne.KeyT: "T",
		fyne.KeyU: "U",
		fyne.KeyV: "V",
		fyne.KeyW: "W",
		fyne.KeyX: "X",
		fyne.KeyY: "Y",
		fyne.KeyZ: "Z",
	}
	for key, letter := range mappings {
		theLetter := letter
		window.Canvas().AddShortcut(
			&desktop.CustomShortcut{KeyName: key},
			func(shortcut fyne.Shortcut) {
				state.typeLetter(theLetter)
				render(&state, window)
			},
		)
	}
	window.Canvas().AddShortcut(
		&desktop.CustomShortcut{KeyName: fyne.KeyBackspace},
		func(shortcut fyne.Shortcut) {
			state.backspace()
			render(&state, window)
		},
	)
	window.Canvas().AddShortcut(
		&desktop.CustomShortcut{KeyName: fyne.KeySpace},
		func(shortcut fyne.Shortcut) {
			err := state.enter()
			if err != nil {
				displayError(err, &state, window)
			} else {
				render(&state, window)
			}
		},
	)

	window.ShowAndRun()
}

var errorTicker *time.Ticker

func displayError(err error, state *AppState, window fyne.Window) {
	if errorTicker != nil {
		errorTicker.Stop()
	}
	state.errorMessage = cases.Title(language.English).String(err.Error())
	errorTicker = time.NewTicker(1 * time.Second)
	render(state, window)

	go func() {
		<-errorTicker.C
		errorTicker.Stop()
		state.errorMessage = ""
		render(state, window)
	}()
}

func render(state *AppState, window fyne.Window) {
	container := container.New(layout.NewVBoxLayout())

	container.Add(header())
	container.Add(statusMessage(state))
	container.Add(wordRows(state))

	space := canvas.NewRectangle(c.Transparent)
	space.SetMinSize(fyne.NewSize(0, 40))
	container.Add(space)

	container.Add(keyboard(state, window))

	window.SetContent(container)
}

func header() *fyne.Container {
	header := container.New(layout.NewVBoxLayout())

	titleRow := container.New(layout.NewMaxLayout())
	iconRect := canvas.NewRectangle(c.Transparent)
	iconRect.SetMinSize(fyne.NewSize(48, 48))
	appIcon, _ := fyne.LoadResourceFromPath("assets/AppIcon.svg")
	icon := widget.NewIcon(appIcon)
	titleRow.Add(container.New(layout.NewHBoxLayout(), container.New(layout.NewMaxLayout(), iconRect, icon)))

	title := canvas.NewText("Gordle", c.Black)
	title.Alignment = fyne.TextAlignCenter
	title.TextSize = 24
	title.TextStyle.Bold = true
	titleRow.Add(title)

	header.Add(titleRow)

	border := canvas.NewRectangle(c.RGBA{R: 211, G: 215, B: 218, A: 255})
	border.SetMinSize(fyne.NewSize(0, 5))
	header.Add(border)

	return header
}

func statusMessage(state *AppState) *fyne.Container {
	message := ""
	if state.game.State == Won {
		message = "You won!"
	} else if state.game.State == Lost {
		message = "The solution was " + state.game.Solution.Letters + "."
	} else if state.errorMessage != "" {
		message = state.errorMessage
	}

	container := container.New(layout.NewMaxLayout())

	blackBox := canvas.NewRectangle(c.Transparent)
	if message != "" {
		blackBox.FillColor = c.Black
	}
	blackBox.SetMinSize(fyne.NewSize(0, 50))
	blackBox.StrokeWidth = 15
	blackBox.StrokeColor = c.White
	container.Add(blackBox)

	statusText := canvas.NewText(message, c.White)
	statusText.Alignment = fyne.TextAlignCenter
	statusText.TextSize = 14
	statusText.TextStyle.Bold = true

	container.Add(statusText)

	return container
}

func wordRows(state *AppState) *fyne.Container {
	rows := container.New(layout.NewVBoxLayout())
	for i, guess := range state.guesses {
		feedback := state.feedbacks[i]
		rows.Add(wordRow(guess, feedback))
	}
	remaining := 6 - len(state.guesses)
	if state.game.State == Started {
		remaining = remaining - 1
		rows.Add(currentWordRow(state.currentWord))
	}
	for i := 0; i < remaining; i++ {
		rows.Add(emptyWordRow())
	}
	return container.New(layout.NewCenterLayout(), rows)
}

func wordRow(word Word, feedback Feedback) *fyne.Container {
	grid := container.New(layout.NewGridLayout(5))
	for i, letter := range strings.Split(word.Letters, "") {
		color := feedback.Colors[i]
		grid.Add(letterBox(letter, color))
	}
	return grid
}

func currentWordRow(word string) *fyne.Container {
	grid := container.New(layout.NewGridLayout(5))
	for _, letter := range strings.Split(word, "") {
		grid.Add(currentWordLetterBox(letter))
	}
	for i := len(word); i < 5; i++ {
		grid.Add(emptyLetterBox())
	}
	return grid
}

func emptyWordRow() *fyne.Container {
	grid := container.New(layout.NewGridLayout(5))
	for i := 0; i < 5; i++ {
		grid.Add(emptyLetterBox())
	}
	return grid
}

func letterBox(letter string, color Color) *fyne.Container {
	var fill c.Color
	switch color {
	case Green:
		fill = c.RGBA{R: 106, G: 170, B: 100, A: 255}
	case Yellow:
		fill = c.RGBA{R: 201, G: 180, B: 88, A: 255}
	case Grey:
		fill = c.RGBA{R: 120, G: 124, B: 126, A: 255}
	}

	box := canvas.NewRectangle(fill)
	box.SetMinSize(fyne.NewSize(62, 62))

	text := canvas.NewText(letter, c.White)
	text.Alignment = fyne.TextAlignCenter
	text.TextSize = 32
	text.TextStyle.Bold = true
	content := container.New(layout.NewMaxLayout(), box, text)

	return content
}

func currentWordLetterBox(letter string) *fyne.Container {
	box := canvas.NewRectangle(c.White)
	box.SetMinSize(fyne.NewSize(62, 62))
	box.StrokeColor = c.RGBA{R: 135, G: 138, B: 140, A: 255}
	box.StrokeWidth = 2.0

	text := canvas.NewText(letter, c.Black)
	text.Alignment = fyne.TextAlignCenter
	text.TextSize = 32
	text.TextStyle.Bold = true
	content := container.New(layout.NewMaxLayout(), box, text)

	return content
}

func emptyLetterBox() *canvas.Rectangle {
	box := canvas.NewRectangle(c.White)
	box.SetMinSize(fyne.NewSize(62, 62))
	box.StrokeColor = c.RGBA{R: 211, G: 214, B: 218, A: 255}
	box.StrokeWidth = 2.0
	return box
}

func keyboard(state *AppState, window fyne.Window) *fyne.Container {
	letterRows := [][]string{
		{"Q", "W", "E", "R", "T", "Y", "U", "I", "O", "P"},
		{"A", "S", "D", "F", "G", "H", "J", "K", "L"},
		{"Z", "X", "C", "V", "B", "N", "M"},
	}
	keyboard := container.New(layout.NewVBoxLayout())
	for i, letters := range letterRows {
		row := container.New(layout.NewHBoxLayout())
		if i == 2 {
			theState := state
			enterButton := widget.NewButton("ENTER", func() {
				err := state.enter()
				if err != nil {
					displayError(err, theState, window)
				} else {
					render(theState, window)
				}
			})
			decoratedEnterButton := decorateButton(enterButton, nil, nil, fyne.NewSize(65.4, 58))
			row.Add(decoratedEnterButton)
		}
		for _, letter := range letters {
			theLetter := letter
			theState := state
			button := widget.NewButton(theLetter, func() {
				theState.typeLetter(theLetter)
				render(theState, window)
			})
			decoratedButton := decorateLetterButton(button, state)
			row.Add(decoratedButton)
		}
		if i == 2 {
			theState := state
			backButton := widget.NewButton("BACK", func() {
				state.backspace()
				render(theState, window)
			})
			decoratedBackButton := decorateButton(backButton, nil, nil, fyne.NewSize(65.4, 58))
			row.Add(decoratedBackButton)
		}
		keyboard.Add(container.New(layout.NewCenterLayout(), row))
	}
	return keyboard
}

func decorateLetterButton(button *widget.Button, state *AppState) *fyne.Container {
	letter := button.Text

	var bgColor c.Color
	var fgColor c.Color
	color, exists := state.letterColors[letter]
	if exists {
		fgColor = c.White
		switch color {
		case Green:
			bgColor = c.RGBA{R: 106, G: 170, B: 100, A: 255}
		case Yellow:
			bgColor = c.RGBA{R: 201, G: 180, B: 88, A: 255}
		case Grey:
			bgColor = c.RGBA{R: 120, G: 124, B: 126, A: 255}
		}
	}

	return decorateButton(button, bgColor, fgColor, fyne.NewSize(43, 58))
}

func decorateButton(button *widget.Button, bgColor c.Color, fgColor c.Color, size fyne.Size) *fyne.Container {
	if bgColor == nil {
		bgColor = c.RGBA{R: 211, G: 214, B: 218, A: 255}
	}
	if fgColor == nil {
		fgColor = c.Black
	}
	letter := button.Text
	button.Text = ""
	rectangle := canvas.NewRectangle(bgColor)
	rectangle.SetMinSize(size)
	text := canvas.NewText(letter, fgColor)
	text.Alignment = fyne.TextAlignCenter
	text.TextStyle.Bold = true
	return container.New(
		layout.NewMaxLayout(),
		rectangle,
		text,
		button,
	)
}
