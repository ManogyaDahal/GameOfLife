package main

import (
	"log"

	"fmt"
	"image/color"
	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"math"
)


const (
	minZoom = 1.0
	maxZoom = 10.0
	zoomSpeed = 0.1 // Controls zoom sensitivity
)

type Game struct {
	World          *World
	worldImg       *ebiten.Image
	pixels         []byte
	paused         bool
	generation     int
	updateInterval int
	camX, camY     float64
	camScale       float64
	lastMouseX     int
	lastMouseY     int
	dragging       bool
	wheelAccum     float64 // Accumulator for smooth wheel scrolling
}

func (g *Game) applyZoom(wheelY float64, mouseX, mouseY int) {
	if wheelY == 0 {
		return
	}

	// Normalize wheel input - clamp to reasonable range
	// Different devices report wildly different values
	normalizedWheel := wheelY
	if normalizedWheel > 3 {
		normalizedWheel = 3
	}
	if normalizedWheel < -3 {
		normalizedWheel = -3
	}

	oldScale := g.camScale
	// Use smaller multiplier for smoother zooming
	newScale := g.camScale * math.Pow(1.0+zoomSpeed, normalizedWheel)

	// Clamp zoom
	if newScale < minZoom {
		newScale = minZoom
	}
	if newScale > maxZoom {
		newScale = maxZoom
	}

	// If scale didn't actually change (hit limits), don't adjust camera
	if newScale == oldScale {
		return
	}

	// Get screen dimensions
	sw, sh := ebiten.WindowSize()
	
	// Ensure mouse is within screen bounds
	mx := float64(mouseX)
	my := float64(mouseY)
	if mx < 0 {
		mx = 0
	}
	if mx > float64(sw) {
		mx = float64(sw)
	}
	if my < 0 {
		my = 0
	}
	if my > float64(sh) {
		my = float64(sh)
	}

	// Calculate the world point under the mouse before zoom
	worldX := (mx - g.camX) / oldScale
	worldY := (my - g.camY) / oldScale

	// Update scale
	g.camScale = newScale

	// Recalculate camera position to keep the same world point under mouse
	g.camX = mx - worldX*g.camScale
	g.camY = my - worldY*g.camScale
}

func (g *Game) clampCamera() {
	sw, sh := ebiten.WindowSize()

	worldW := float64(g.World.width) * g.camScale
	worldH := float64(g.World.height) * g.camScale

	// If world smaller than screen, center it
	if worldW <= float64(sw) {
		g.camX = (float64(sw) - worldW) / 2
	} else {
		// Allow world to be dragged until edge reaches screen edge
		maxX := 0.0
		minX := float64(sw) - worldW
		if g.camX > maxX {
			g.camX = maxX
		}
		if g.camX < minX {
			g.camX = minX
		}
	}

	if worldH <= float64(sh) {
		g.camY = (float64(sh) - worldH) / 2
	} else {
		maxY := 0.0
		minY := float64(sh) - worldH
		if g.camY > maxY {
			g.camY = maxY
		}
		if g.camY < minY {
			g.camY = minY
		}
	}
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
	// Fill screen with black to prevent artifacts
	screen.Fill(color.RGBA{0, 0, 0, 255})
	
	g.World.Draw(g.pixels)
	g.worldImg.WritePixels(g.pixels)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(g.camScale, g.camScale)
	op.GeoM.Translate(g.camX, g.camY)
	screen.DrawImage(g.worldImg, op)

	if g.paused {
		ebitenutil.DebugPrint(screen,
		fmt.Sprintf("PAUSED | Zoom: %.2f \nInstructions: space-play c-clear r-random", g.camScale))
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *Game) interact() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		worldX := int((float64(x) - g.camX) / g.camScale)
		worldY := int((float64(y) - g.camY) / g.camScale)

		if worldX >= 0 && worldX < g.World.width &&
		worldY >= 0 && worldY < g.World.height {
			g.World.area[worldY*g.World.width+worldX] = true
		}
	}
}

func (g *Game) handleCamera() {
	curX, curY := ebiten.CursorPosition()

	// -------- ZOOM (handle first, before any clamping) --------
	_, wheelY := ebiten.Wheel()
	if wheelY != 0 {
		// Accumulate wheel delta for smoother experience
		g.wheelAccum += wheelY
		
		// Process accumulated wheel movement in small increments
		const wheelThreshold = 0.3 // Adjust this for sensitivity
		
		if math.Abs(g.wheelAccum) >= wheelThreshold {
			// Calculate how many steps to zoom
			steps := math.Floor(g.wheelAccum / wheelThreshold)
			g.applyZoom(steps, curX, curY)
			
			// Remove processed amount from accumulator
			g.wheelAccum -= steps * wheelThreshold
			
			// Clamp immediately after zoom to prevent black areas
			g.clampCamera()
		}
	}

	// -------- PANNING (Right Mouse) --------
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		if !g.dragging {
			g.dragging = true
			g.lastMouseX = curX
			g.lastMouseY = curY
		} else {
			dx := float64(curX - g.lastMouseX)
			dy := float64(curY - g.lastMouseY)

			g.camX += dx
			g.camY += dy

			g.lastMouseX = curX
			g.lastMouseY = curY
			
			// Clamp after panning
			g.clampCamera()
		}
	} else {
		g.dragging = false
	}
}

func main() {

	x, y := ebiten.Monitor().Size()
	g := &Game{
		World:    NewWorld(x, y, int((x*y)/10)),
		camScale: 1,
		paused:   true,
	}

  g.worldImg = ebiten.NewImage(g.World.width, g.World.height)
  g.pixels = make([]byte, g.World.width*g.World.height*4)

	ebiten.SetFullscreen(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Game of life")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
