package main

import (
	"log"

	"fmt"
	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"math"
	// "golang.org/x/text/width"
)

// const (
// 	screenWidth  = 320
// 	screenHeight = 240
// )

type Game struct {
	World          *World
	pixels         []byte
	paused         bool
	generation     int
	updateInterval int
	camX, camY     float64
	camScale       float64
	lastMouseX     int
	lastMouseY     int
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

	g.handleCamera()
	g.interact()
	if !g.paused {
		g.World.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	worldImage := ebiten.NewImage(g.World.width, g.World.height)
	if g.pixels == nil {
		g.pixels = make([]byte, g.World.width*g.World.height*4)
	}
	g.World.Draw(g.pixels)
	worldImage.WritePixels(g.pixels)
	if g.paused {
		ebitenutil.DebugPrint(screen, "PAUSED")
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(g.camScale, g.camScale)
	op.GeoM.Translate(g.camX, g.camY)
	screen.DrawImage(worldImage, op)

	if g.paused {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("PAUSED | Zoom: %.2f", g.camScale))
	}
}

// func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
// 	return screenWidth, screenHeight
// }

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.World.width = outsideWidth
	g.World.height = outsideHeight

	return outsideWidth, outsideHeight
}

// Sets the cell at the position that the cursor was left clicked to alive
func (g *Game) interact() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		worldX := int((float64(x) - g.camX) / g.camScale)
		worldY := int((float64(y) - g.camY) / g.camScale)

		// if worldX >= 0 && worldX < g.World.width &&
		// 	worldY >= 0 && worldY < g.World.height {
		g.World.area[worldY*g.World.width+worldX] = true
	}
	// }
}

func (g *Game) handleCamera() {
	curX, curY := ebiten.CursorPosition()

	// PANNING: Right Mouse Button
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		g.camX += float64(curX - g.lastMouseX)
		g.camY += float64(curY - g.lastMouseY)
	}
	g.lastMouseX, g.lastMouseY = curX, curY

	// ZOOMING: Mouse Wheel
	_, wheelY := ebiten.Wheel()
	if wheelY != 0 {
		oldScale := g.camScale

		// Touchpads move in small increments.
		// Using math.Pow makes the zoom feel "smooth" and consistent.
		// 0.05 is the sensitivity. Increase it if it's too slow.
		g.camScale *= math.Pow(1.1, wheelY)

		// Safety bounds
		if g.camScale < 0.1 {
			g.camScale = 0.1
		}
		if g.camScale > 20.0 {
			g.camScale = 20.0
		}

		// Zoom toward cursor (Keep the pixel under the cursor stationary)
		g.camX -= float64(curX) * (g.camScale/oldScale - 1)
		g.camY -= float64(curY) * (g.camScale/oldScale - 1)
	}
}

func main() {

	x, y := ebiten.Monitor().Size()
	g := &Game{
		World:    NewWorld(x, y, int((x*y)/10)),
		camScale: 1.0,
		paused:   true,
	}

	ebiten.SetFullscreen(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Game of life")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
