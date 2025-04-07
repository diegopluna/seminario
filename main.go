package main

import (
	"fmt"
	"math"
	"seminario/astar"
)

type GridPoint struct {
	X, Y int
}

type Grid struct {
	Width, Height int
	Obstacles     map[GridPoint]bool
}

func NewGrid(width, height int) *Grid {
	return &Grid{
		Width:     width,
		Height:    height,
		Obstacles: make(map[GridPoint]bool),
	}
}

func (g *Grid) AddObstacle(x, y int) {
	g.Obstacles[GridPoint{x, y}] = true
}

func (g *Grid) Neighbors(node astar.Node) []astar.Node {
	p := node.(GridPoint)
	neighbors := []astar.Node{}
	moves := []GridPoint{
		{p.X + 1, p.Y},
		{p.X - 1, p.Y},
		{p.X, p.Y + 1},
		{p.X, p.Y - 1},
	}

	for _, move := range moves {
		if move.X >= 0 && move.X < g.Width && move.Y >= 0 && move.Y < g.Height {
			if !g.Obstacles[move] {
				neighbors = append(neighbors, move)
			}
		}
	}

	return neighbors
}

func (g *Grid) Cost(from, to astar.Node) float64 {
	return 1.0
}

func ManhattanDistance(a, b astar.Node) float64 {
	p1 := a.(GridPoint)
	p2 := b.(GridPoint)
	return math.Abs(float64(p1.X-p2.X)) + math.Abs(float64(p1.Y-p2.Y))
}

func main() {
	grid := NewGrid(10, 10)

	grid.AddObstacle(1, 1)
	grid.AddObstacle(1, 2)
	grid.AddObstacle(2, 2)
	grid.AddObstacle(3, 2)
	grid.AddObstacle(4, 2)
	grid.AddObstacle(5, 2)
	grid.AddObstacle(5, 1)
	grid.AddObstacle(5, 0)

	start := GridPoint{X: 0, Y: 0}
	goal := GridPoint{X: 6, Y: 0}

	fmt.Printf("Finding path from %v to %v\n", start, goal)

	path, cost, err := astar.AStar(start, goal, grid, ManhattanDistance)

	if err != nil {
		fmt.Printf("Error finding path: %v\n", err)
	} else {
		fmt.Printf("Path found with cost %.2f:\n", cost)
		for i, p := range path {
			gp := p.(GridPoint) // Type assertion
			fmt.Printf(" %d: (%d, %d)\n", i, gp.X, gp.Y)
		}
		printGridWithPath(grid, path, start, goal)
	}

	fmt.Println("\nTrying a case with no possible path:")
	goalUnreachable := GridPoint{X: 1, Y: 1} // Inside an obstacle block
	_, _, err = astar.AStar(start, goalUnreachable, grid, ManhattanDistance)
	if err != nil {
		fmt.Printf("Correctly failed: %v\n", err)
	} else {
		fmt.Println("Error: Found a path where none should exist.")
	}
}

func printGridWithPath(
	grid *Grid,
	path []astar.Node,
	start, goal GridPoint,
) {
	pathSet := make(map[GridPoint]bool)
	for _, p := range path {
		pathSet[p.(GridPoint)] = true
	}

	fmt.Println("\nGrid Visualization:")
	fmt.Print("   ")
	for x := range grid.Width {
		fmt.Printf("%2d", x)
	}
	fmt.Println()
	fmt.Print("  +-")
	for range grid.Width {
		fmt.Print("--")
	}
	fmt.Println("-+")

	for y := range grid.Height {
		fmt.Printf("%2d| ", y)
		for x := range grid.Width {
			p := GridPoint{x, y}
			if p == start {
				fmt.Print("S ")
			} else if p == goal {
				fmt.Print("G ")
			} else if grid.Obstacles[p] {
				fmt.Print("##") // Obstacle
			} else if pathSet[p] {
				fmt.Print("..") // Path
			} else {
				fmt.Print("  ") // Empty
			}
		}
		fmt.Println(" |")
	}
	fmt.Print("  +-")
	for range grid.Width {
		fmt.Print("--")
	}
	fmt.Println("-+")
}
