package main

import (
	gordle "gordle/src"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	gordle.StartFyneGame()
}
