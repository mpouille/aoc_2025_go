package main

import (
	"github.com/jpillora/puzzler/harness/aoc"
	"strings"
	"fmt"
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

	lines := strings.Split(input, "\n")
	grid := make([][]rune, len(lines))
	startingPosition := []int{0, 0}
	for i, line := range lines {
		grid[i] = []rune(line)
		if i == 0 {
			for j, char := range line {
				if char == 'S' {
					startingPosition = []int{i, j}
					break
				}
			}
		}
	}
	// when you're ready to do part 2, remove this "not implemented" block
	if part2 {
		return runPart2(grid, startingPosition)
	}
	// solve part 1 here
	return runPart1(grid, startingPosition)
}

func checkForCollisions(grid [][]rune, i int, j int) bool {
	if grid[i][j] == '^' || grid[i][j] == 'O' {
		return true
	}
	return false
}

func rayGoDown(grid [][]rune, i int, j int, part2 bool, memo map[string]int) int {
	key := fmt.Sprintf("%d,%d", i, j)
	if part2 {
		if val, ok := memo[key]; ok {
			return val
		}
	}

	newI := i + 1
	if (newI == len(grid)) {
		memo[key] = 1
		return 1
	}
	if checkForCollisions(grid, newI, j) {
		if !part2 && grid[newI][j] == 'O' {
			return 1
		}
		grid[newI][j + 1] = '|'
		grid[newI][j - 1] = '|'
		grid[newI][j] = 'O'
		val := rayGoDown(grid, newI, j + 1, part2, memo) + rayGoDown(grid, newI, j - 1, part2, memo)
		memo[key] = val
		return val
	}
	grid[newI][j] = '|'
	val := rayGoDown(grid, newI, j, part2, memo)
	memo[key] = val
	return val
}

func runPart1(grid [][]rune, startingPosition []int) int {
	memo := make(map[string]int)
	return rayGoDown(grid, startingPosition[0], startingPosition[1], false, memo) - 1
}

func runPart2(grid [][]rune, startingPosition []int) int {
	memo := make(map[string]int)
	return rayGoDown(grid, startingPosition[0], startingPosition[1], true, memo)
}

func echo(grid [][]rune) {
	for _, row := range grid {
		for _, char := range row {
			fmt.Print(string(char))
		}
		fmt.Println()
	}
}