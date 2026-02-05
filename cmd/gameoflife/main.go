package main

import (
	"log"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  =  320
	screenHeight = 	240
)

type Game struct { 
	World *World
	pixels []byte
}

// game struct method
func (g *Game)Update() error { 
	g.World.Update()
	return nil
}

func (g *Game)Draw (screen *ebiten.Image) { 
	if g.pixels == nil {
		g.pixels = make([]byte, screenWidth*screenHeight*4)
	}
	g.World.Draw(g.pixels)
	screen.WritePixels(g.pixels)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := &Game{
		World: NewWorld(screenWidth, screenHeight, int((screenWidth*screenHeight)/10)),
	}

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Game of life")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
