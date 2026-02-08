# Game of Life

Minimal Conway's Game of Life implemented in Go using Ebiten.

### Quick start
1. From project root:
- Run: `go mod tidy` to install dependencies
- Run: `go run ./cmd/gameoflife/*`

### Controls
- Space: Pause / unpause
- R (paused): Randomize the board
- C (paused): Clear the board
- Left click: Set a cell alive

### Source
- `cmd/gameoflife/main.go` — app entry, input handling, Ebiten game loop
- `cmd/gameoflife/logic.go` — world state, update rules, drawing

### File Tree 
```File Tree 
.
├── cmd
│   ├── gameoflife     
│   │   ├── logic.go
│   │   └── main.go
│   └── server
│       └── server.go
├── go.mod
├── go.sum
├── LICENSE
├── README.md
└── static                
```
### License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
