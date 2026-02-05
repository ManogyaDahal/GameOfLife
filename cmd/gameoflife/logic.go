package main

import (
	"math/rand"
)

// World represents the game state
type World struct{ 
	area []bool
	width int 
	height int
}

// New world creates a new world
func NewWorld(width int, height int, maxInitLiveCells int) *World { 
  w := &World{
		area: make([]bool, width*height),
		width: width, 
		height: height,
	}

	return w
}

// Initializes the world with random no. of cells
func (w *World) init(maxInitLiveCells int) {
	for i := 0; i < maxInitLiveCells; i++	{ 
		x := rand.Intn(w.width)
		y := rand.Intn(w.height)

		w.area[y*w.width + x] = true
	}
}

//update the game state by one trick
func (w *World) Update() {
	width := w.width
	height := w.height
	next := make([]bool, width*height)

	for y:=0; y < height; y++ {
		for x := 0; x < width; x++ { 
			pop :=	w.neighbourCount(x, y) 
			switch { 
			case pop < 2:
				// rule 1. Any live cell with fewer than two live neighbours
				// dies, as if caused by under-population.
				next[y*width+x] = false

			case (pop == 2 || pop == 3) && w.area[y*width+x]:
				// rule 2. Any live cell with two or three live neighbours
				// lives on to the next generation.
				next[y*width+x] = true

			case pop > 3:
				// rule 3. Any live cell with more than three live neighbours
				// dies, as if by over-population.
				next[y*width+x] = false

			case pop == 3:
				// rule 4. Any dead cell with exactly three live neighbours
				// becomes a live cell, as if by reproduction.
				next[y*width+x] = true
			}
		}
	}
	w.area = next
}

// returns the number of neighbours
func (w *World)neighbourCount (x int, y int) int { 
	c := 0 
	for j := -1; j <= 1; j++{ 
		for i := -1; i <= 1; i++{
			if i == 0 && j == 0 { continue }

			x2 := x + i
			y2 := y + j
			if x2 < 0 || y2 < 0 || w.width <= x2 || w.height <= y2 { continue }

			if w.area[y2*w.width +x2]{ 
				c++
			}
		}
	}
	return c
}

//draw paint in current game state
func (w *World) Draw(pix []byte) { 
	for i, v := range w.area{
		if v {
			pix[4*i] = 0xff
			pix[4*i+1] = 0xff
			pix[4*i+2] = 0xff
			pix[4*i+3] = 0xff
		} else {
			pix[4*i] = 0
			pix[4*i+1] = 0
			pix[4*i+2] = 0
			pix[4*i+3] = 0
		}
	}
}
