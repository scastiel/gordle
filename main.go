package main

import (
	gordle "gordle/src"
	"math/rand"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	os.Setenv("FYNE_THEME", "light")
	gordle.StartFyneGame()
}
