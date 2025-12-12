package main

import (
	"github.com/jpillora/puzzler/harness/aoc"
	"strings"
)

func main() {
	aoc.Harness(run)
}

// on code change, run will be executed 4 times:
// 1. with: false (part1), and example input
// 2. with: true (part2), and example input
// 3. with: false (part1), and user input
// 4. with: true (part2), and user input
// the return value of each run is printed to stdout
func run(part2 bool, input string) any {
	// when you're ready to do part 2, remove this "not implemented" block
	if part2 {
		return partTwoFn(input)
	}
	// solve part 1 here
	return partOneFn(input)
}

func partOneFn(input string) int {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	result := 0

	type Cell struct {
		isRollpaper           bool
		adjacentRollpaperCount int
	}
	
	grid := make([][]Cell, len(lines))
	for i, line := range lines {
		grid[i] = make([]Cell, len(line))
		for j, char := range line {
			grid[i][j] = Cell{
				isRollpaper:           char == '@',
				adjacentRollpaperCount: 0, 
			}
		}
	}
	
	// Count adjacent rollpapers for each cell
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			count := 0
			// Check all 8 adjacent positions
			for di := -1; di <= 1; di++ {
				for dj := -1; dj <= 1; dj++ {
					if di == 0 && dj == 0 {
						continue // Skip the cell itself
					}
					ni, nj := i+di, j+dj
					if ni >= 0 && ni < len(grid) && nj >= 0 && nj < len(grid[ni]) {
						if grid[ni][nj].isRollpaper {
							count++
						}
					}
				}
			}
			grid[i][j].adjacentRollpaperCount = count
			if grid[i][j].adjacentRollpaperCount < 4 && grid[i][j].isRollpaper {
				result++
			}
		}
	}
	
	return result
}



func partTwoFn(input string) int {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	result := 0
	
	type Cell struct {
		isRollpaper           bool
		adjacentRollpaperCount int
	}
	
	grid := make([][]Cell, len(lines))
	for i, line := range lines {
		grid[i] = make([]Cell, len(line))
		for j, char := range line {
			grid[i][j] = Cell{
				isRollpaper:           char == '@',
				adjacentRollpaperCount: 0, 
			}
		}
	}
	

	initialGrid := make([][]Cell, len(grid))

	for {
		for i := range grid {
			initialGrid[i] = make([]Cell, len(grid[i]))
			copy(initialGrid[i], grid[i])
		}
		partialResult := 0
		// Count adjacent rollpapers for each cell
		for i := 0; i < len(initialGrid); i++ {
			for j := 0; j < len(initialGrid[i]); j++ {
				count := 0
				// Check all 8 adjacent positions
				for di := -1; di <= 1; di++ {
					for dj := -1; dj <= 1; dj++ {
						if di == 0 && dj == 0 {
							continue // Skip the cell itself
						}
						ni, nj := i+di, j+dj
						if ni >= 0 && ni < len(initialGrid) && nj >= 0 && nj < len(initialGrid[ni]) {
							if initialGrid[ni][nj].isRollpaper {
								count++
							}
						}
					}
				}
				initialGrid[i][j].adjacentRollpaperCount = count
				if initialGrid[i][j].adjacentRollpaperCount < 4 && initialGrid[i][j].isRollpaper {
					partialResult++
					grid[i][j] = Cell{
						isRollpaper: false,
						adjacentRollpaperCount: 0, 
					}
				}
			}
		}
		result += partialResult

		if partialResult == 0 {
			break
		}
	}
	
	return result
}