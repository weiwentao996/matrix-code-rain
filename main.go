package main

import (
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	density      = 0.01 // 字符密度[0,1]
	screenWidth  = 600
	screenHeight = 400
	minXSpan     = 2
	minYSpan     = 30
	maxYLength   = 30
	minYLength   = 4
	fullChar     = 0
)

type Point struct {
	Char        rune
	Color       termbox.Attribute
	YLineIndex  int
	YLineLength int
}

var globalMap [][]Point

var emptyPoint = Point{
	Char:  fullChar,
	Color: termbox.ColorGreen,
}

func init() {
	globalMap = make([][]Point, screenHeight, screenHeight)
	for y, _ := range globalMap {
		for x, _ := range globalMap[y] {
			globalMap[y][x] = emptyPoint
		}
	}

}

var charPool = []rune{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o',
	'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O',
	'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	';', '.', '/', '*', '!', '@', '#', '$', '%', '^', '&', '(', ')', '_', '+',
	'-', '=', '`', '~',
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	rand.Seed(time.Now().UnixNano())
	for i, _ := range globalMap {
		globalMap[i] = make([]Point, screenWidth)
	}
	for {
		g := globalMap
		_ = g
		// Clear the screen
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

		// Copy previous line
		for y := len(globalMap) - 1; y > 0; y-- {
			globalMap[y] = globalMap[y-1]
		}

		// Clear the first line
		globalMap[0] = make([]Point, screenWidth)
		for i, _ := range globalMap[0] {
			globalMap[0][i] = Point{
				Char: fullChar,
			}
		}

		// Draw the line
		for x := 0; x < screenWidth; x++ {
			if globalMap[1][x].Char != fullChar {
				if globalMap[1][x].YLineIndex < globalMap[1][x].YLineLength {
					globalMap[0][x] = Point{
						Char:        charPool[rand.Intn(len(charPool))],
						Color:       termbox.ColorGreen,
						YLineLength: globalMap[1][x].YLineLength,
						YLineIndex:  globalMap[1][x].YLineIndex + 1,
					}
				} else {
					globalMap[0][x] = emptyPoint
				}
			}
		}

		// Draw the line numbers
		for i := 0; i < screenWidth*density; i++ {
		nextChar:
			x := rand.Intn(screenWidth)
			if globalMap[1][x].Char == fullChar {
				for y := 1; y <= minYSpan; y++ {
					if globalMap[y][x].Char != fullChar {
						goto nextChar
					}
				}

				// Judge whether there are characters near the position
				for i := 1; i <= minXSpan; i++ {
					if x-i >= 0 {
						if globalMap[0][x-i].Char != fullChar {
							goto nextChar
						}
					}
					if x+i < screenWidth {
						if globalMap[0][x+i].Char != fullChar {
							goto nextChar
						}
					}
				}
				globalMap[0][x] = Point{
					Char:        charPool[rand.Intn(len(charPool))],
					Color:       termbox.ColorDefault,
					YLineIndex:  0,
					YLineLength: rand.Intn(maxYLength-minYLength) + minYLength,
				}
			}
		}

		for y, _ := range globalMap {
			for x, _ := range globalMap[y] {
				if globalMap[y][x].Char != fullChar {
					termbox.SetCell(x, y, globalMap[y][x].Char, globalMap[y][x].Color, termbox.ColorDefault)
				}
			}
		}

		// Update the screen
		termbox.Flush()

		// Wait for a short amount of time
		time.Sleep(time.Millisecond * 50)
	}
}
