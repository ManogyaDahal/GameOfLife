package main

import (
	"log"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	// "golang.org/x/text/width"
)

// const (
// 	screenWidth  = 320
// 	screenHeight = 240
// )

type Game struct {
	width          int
	height         int
	World          *World
	pixels         []byte
	paused         bool
	generation     int
	updateInterval int
}

// game struct method
func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.paused = !g.paused
	}
	if g.paused && inpututil.IsKeyJustPressed(ebiten.KeyC) {
		g.World.area = make([]bool, g.World.width*g.World.height)
	}
	if g.paused && inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.World.init((g.World.width * g.World.height) / 10) //max number of live cells
	}

	g.interact()
	if !g.paused {
		g.World.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.pixels == nil {
		g.pixels = make([]byte, g.World.width*g.World.height*4)
	}
	g.World.Draw(g.pixels)
	screen.WritePixels(g.pixels)
	if g.paused {
		ebitenutil.DebugPrint(screen, "PAUSED")
	}
}

// func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
// 	return screenWidth, screenHeight
// }

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.width = outsideWidth
	g.height = outsideHeight

	return outsideWidth, outsideHeight
}

// Sets the cell at the position that the cursor was left clicked to alive
func (g *Game) interact() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		if x >= 0 && x < g.World.width &&
			y >= 0 && y < g.World.height {

			g.World.area[y*g.World.width+x] = true
		}
	}
}

func main() {

	x, y := ebiten.Monitor().Size()
	g := &Game{
		World:  NewWorld(x, y, int((x*y)/10)),
		paused: true,
	}

	ebiten.SetFullscreen(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Game of life")
	ebiten.SetTPS(12)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
