package main

import (
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	screenWidth  = 1000
	screenHeight = 1000
	minXSpan     = 3
	fullChar     = rune(' ')
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	rand.Seed(time.Now().UnixNano())
	globalMap := make([][]int, screenHeight)
	for i, _ := range globalMap {
		globalMap[i] = make([]int, screenWidth)
	}
	for {
		// Clear the screen
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		for y := len(globalMap) - 1; y > 0; y-- {
			globalMap[y] = globalMap[y-1]
		}
		globalMap[0] = make([]int, screenWidth)
		for i, _ := range globalMap[0] {
			globalMap[0][i] = int(fullChar)
		}

		// Draw the numbers
		for i := 0; i < screenWidth/2; i++ {
		nextChar:
			x := rand.Intn(screenWidth)
			// Judge whether there are characters near the position
			for i := 1; i <= minXSpan; i++ {
				if x-i >= 0 {
					if globalMap[0][x-i] != int(fullChar) {
						goto nextChar
					}
				}
				if x+i < screenWidth {
					if globalMap[0][x+i] != int(fullChar) {
						goto nextChar
					}
				}
			}
			ch := rand.Intn(2)
			chr := '0' + rune(ch)
			globalMap[0][x] = int(chr)
		}
		for y, _ := range globalMap {
			for x, _ := range globalMap[y] {
				if globalMap[y][x] != 0 {
					termbox.SetCell(x, y, rune(globalMap[y][x]), termbox.ColorGreen, termbox.ColorDefault)
				}
			}
		}

		// Update the screen
		termbox.Flush()

		// Wait for a short amount of time
		time.Sleep(time.Millisecond * 50)
	}
}
