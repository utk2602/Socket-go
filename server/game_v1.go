package main

import (
	"fmt"
	"math/rand"
)

type Game struct {
	Grid    [][]string
	Players map[string]*Player
	Size    int
}

type Player struct {
	Name string
	X, Y int
	Alive bool
}

func NewGame(size int) *Game {
	grid := make([][]string, size)
	for i := range grid {
		grid[i] = make([]string, size)
	}

	game := &Game{
		Grid:    grid,
		Players: make(map[string]*Player),
		Size:    size,
	}

	// place mines
	for i := 0; i < size; i++ {
		x, y := rand.Intn(size), rand.Intn(size)
		grid[x][y] = "M"
	}

	return game
}

func (g *Game) AddPlayer(name string) {
	if len(g.Players) >= 2 {
		return
	}
	startX, startY := 0, len(g.Players)
	g.Players[name] = &Player{Name: name, X: startX, Y: startY, Alive: true}
}

func (g *Game) MovePlayer(name, direction string) string {
	player, ok := g.Players[name]
	if !ok {
		g.AddPlayer(name)
		player = g.Players[name]
	}

	if !player.Alive {
		return fmt.Sprintf("%s is out!", name)
	}

	oldX, oldY := player.X, player.Y
	switch direction {
	case "UP":
		if player.X > 0 {
			player.X--
		}
	case "DOWN":
		if player.X < g.Size-1 {
			player.X++
		}
	case "LEFT":
		if player.Y > 0 {
			player.Y--
		}
	case "RIGHT":
		if player.Y < g.Size-1 {
			player.Y++
		}
	default:
		return "Invalid move."
	}

	if g.Grid[player.X][player.Y] == "M" {
		player.Alive = false
		return fmt.Sprintf("💥 %s stepped on a mine at (%d,%d)!", name, player.X, player.Y)
	}

	if player.X == g.Size-1 {
		return fmt.Sprintf("🏁 %s reached the end and wins!", name)
	}

	return fmt.Sprintf("%s moved from (%d,%d) → (%d,%d)", name, oldX, oldY, player.X, player.Y)
}

func (g *Game) Display(requester string) string {
	output := "\nCurrent Board:\n"
	for i := 0; i < g.Size; i++ {
		for j := 0; j < g.Size; j++ {
			cell := "."
			for _, p := range g.Players {
				if p.X == i && p.Y == j && p.Alive {
					cell = string(p.Name[len(p.Name)-1]) // show player number
				}
			}
			output += cell + " "
		}
		output += "\n"
	}
	return output
}
