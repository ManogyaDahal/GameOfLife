package main

import (
	"log"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  =  320
	screenHeight = 	240
)

type Game struct { 
	World *World
	pixels []byte
	paused  bool
}

// game struct method
func (g *Game)Update() error { 
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.paused = !g.paused
	}
	if g.paused && inpututil.IsKeyJustPressed(ebiten.KeyC) {
		g.World.area = make([]bool, g.World.width*g.World.height)
	}
	if g.paused && inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.World.init((screenHeight*screenWidth)/10) //max number of live cells
	}

	g.interact()
	if !g.paused {
		g.World.Update()
	} 
	return nil
}

func (g *Game)Draw (screen *ebiten.Image) { 
	if g.pixels == nil {
		g.pixels = make([]byte, screenWidth*screenHeight*4)
	}
	g.World.Draw(g.pixels)
	screen.WritePixels(g.pixels)
	if g.paused { ebitenutil.DebugPrint(screen, "PAUSED") }
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// Sets the cell at the position that the cursor was left clicked to alive
func (g *Game) interact() {
    if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
        x, y := ebiten.CursorPosition()

        if x >= 0 && x < g.World.width &&
           y >= 0 && y < g.World.height {

            g.World.area[y*g.World.width + x] = true
        }
    }
}

func main() {
	g := &Game{
		World: NewWorld(screenWidth, screenHeight, int((screenWidth*screenHeight)/10)),
		paused: true,
	}

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Game of life")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
