package main

import (
	"math/rand"
)

// World represents the game state
type World struct {
	area   []bool
	age    []int
	width  int
	height int
}

// New world creates a new world
func NewWorld(width int, height int, maxInitLiveCells int) *World {
	w := &World{
		area:   make([]bool, width*height),
		age:    make([]int, width*height),
		width:  width,
		height: height,
	}

	return w
}

// Initializes the world with random no. of cells
func (w *World) init(maxInitLiveCells int) {
	w.age = make([]int, w.width*w.height)
	for i := 0; i < maxInitLiveCells; i++ {
		x := rand.Intn(w.width)
		y := rand.Intn(w.height)

		w.area[y*w.width+x] = true
		w.age[y*w.width+x] = 1
	}
}

// update the game state by one trick
func (w *World) Update() {
	width := w.width
	height := w.height
	next := make([]bool, width*height)
	nextAge := make([]int, width*height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := y*width + x
			pop := w.neighbourCount(x, y)
			switch {
			case pop < 2:
				// rule 1. Any live cell with fewer than two live neighbours
				// dies, as if caused by under-population.
				next[idx] = false

			case (pop == 2 || pop == 3) && w.area[idx]:
				// rule 2. Any live cell with two or three live neighbours
				// lives on to the next generation.
				next[idx] = true
				nextAge[idx] = w.age[idx] + 1

			case pop > 3:
				// rule 3. Any live cell with more than three live neighbours
				// dies, as if by over-population.
				next[idx] = false

			case pop == 3:
				// rule 4. Any dead cell with exactly three live neighbours
				// becomes a live cell, as if by reproduction.
				next[idx] = true
				nextAge[idx] = 1
			}
		}
	}
	w.area = next
	w.age = nextAge
}

// returns the number of neighbours
func (w *World) neighbourCount(x int, y int) int {
	c := 0
	for j := -1; j <= 1; j++ {
		for i := -1; i <= 1; i++ {
			if i == 0 && j == 0 {
				continue
			}

			x2 := x + i
			y2 := y + j
			if x2 < 0 || y2 < 0 || w.width <= x2 || w.height <= y2 {
				continue
			}

			if w.area[y2*w.width+x2] {
				c++
			}
		}
	}
	return c
}

// draw paint in current game state
// Cells are coloured by age: green (newborn) → yellow → white (old)
func (w *World) Draw(pix []byte) {
	for i, v := range w.area {
		if v {
			a := w.age[i]
			var r, g, b byte
			switch {
			case a <= 1:
				r, g, b = 0x00, 0xff, 0x00 // green
			case a <= 5:
				r, g, b = 0x80, 0xff, 0x00 // chartreuse
			case a <= 15:
				r, g, b = 0xff, 0xff, 0x00 // yellow
			case a <= 30:
				r, g, b = 0xff, 0xcc, 0x00 // amber
			default:
				r, g, b = 0xff, 0xff, 0xff // white
			}
			pix[4*i] = r
			pix[4*i+1] = g
			pix[4*i+2] = b
			pix[4*i+3] = 0xff
		} else {
			pix[4*i] = 0
			pix[4*i+1] = 0
			pix[4*i+2] = 0
			pix[4*i+3] = 0
		}
	}
}
