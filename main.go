package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	density      = 0.01 // 字符密度[0,1]
	screenWidth  = 1000
	screenHeight = 1000
	minXSpan     = 6
	minYSpan     = 30
	maxYLength   = 30
	minYLength   = 4
	fullChar     = 11
)

type Point struct {
	Char        rune
	Color       termbox.Attribute
	YLineIndex  int
	YLineLength int
	StringIndex int
}

var globalMap [][]Point
var globalStringIndex int

var emptyPoint = Point{
	Char:  fullChar,
	Color: termbox.ColorGreen,
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

var stringPool []string

func init() {
	globalMap = make([][]Point, screenHeight, screenHeight)
	for y := 0; y < screenHeight; y++ {
		var item []Point
		for x := 0; x < screenWidth; x++ {
			item = append(item, emptyPoint)
		}
		globalMap[y] = item
	}
	file, err := os.OpenFile("./text", os.O_RDONLY, 0444)
	defer file.Close()
	if err != nil {
		return
	}

	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		str := string(line)
		if err != nil && err == io.EOF {
			break
		}
		if strings.TrimSpace(str) != "" && len([]rune(str)) < screenWidth {
			stringPool = append(stringPool, strings.TrimSpace(str))
		}
	}
	s := stringPool
	_ = s
	fmt.Println(s)
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	rand.Seed(time.Now().UnixNano())
	for {
		// Clear the screen
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		extends()
		genFirstLine()
		drawLine()
		// Update the screen
		termbox.Flush()

		// Wait for a short amount of time
		time.Sleep(time.Millisecond * 50)
	}
}

func drawLine() {
	for y, _ := range globalMap {
		for x, _ := range globalMap[y] {
			if globalMap[y][x].Char != fullChar {
				termbox.SetCell(x, y, globalMap[y][x].Char, globalMap[y][x].Color, termbox.ColorDefault)
			}
		}
	}
}

func genFirstLine() {
	// Clear the first line
	globalMap[0] = make([]Point, screenWidth)
	for i, _ := range globalMap[0] {
		globalMap[0][i] = emptyPoint
	}
	var char rune
	for x := 0; x < screenWidth; x++ {
		if globalMap[1][x].Char != fullChar {
			if globalMap[1][x].YLineIndex < globalMap[1][x].YLineLength-1 {
				if stringPool == nil || len(stringPool) == 0 {
					char = charPool[rand.Intn(len(charPool))]
				} else {
					char = []rune(stringPool[globalMap[1][x].StringIndex])[globalMap[1][x].YLineIndex+1]
				}

				globalMap[0][x] = Point{
					Char:        char,
					Color:       termbox.ColorGreen,
					StringIndex: globalMap[1][x].StringIndex,
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
			var length, stringIndex int
			if stringPool == nil || len(stringPool) == 0 {
				char = charPool[rand.Intn(len(charPool))]
				length = rand.Intn(maxYLength-minYLength) + minYLength
			} else {
				globalStringIndex = globalStringIndex % len(stringPool)
				stringIndex = globalStringIndex
				globalStringIndex++
				char = []rune(stringPool[stringIndex])[0]
				length = len([]rune(stringPool[stringIndex]))
			}
			globalMap[0][x] = Point{
				Char:        char,
				Color:       termbox.ColorDefault,
				StringIndex: stringIndex,
				YLineIndex:  0,
				YLineLength: length,
			}
		}
	}

}

func extends() {
	for y := len(globalMap) - 1; y > 0; y-- {
		globalMap[y] = globalMap[y-1]
	}
}
