package main

import (
	"github.com/jpillora/puzzler/harness/aoc"
	"strings"
	"strconv"
	"sort"
)

func main() {
	aoc.Harness(run)
}

type Range struct {
	start int
	end   int
}

// on code change, run will be executed 4 times:
// 1. with: false (part1), and example input
// 2. with: true (part2), and example input
// 3. with: false (part1), and user input
// 4. with: true (part2), and user input
// the return value of each run is printed to stdout
func run(part2 bool, input string) any {
	lines := splitInput(input)
	availableIngredients := make([]Range, len(strings.Split(lines[0], "\n")))
	for i, line := range strings.Split(lines[0], "\n") {
		parts := strings.Split(line, "-")
		start, _ := strconv.Atoi(parts[0])
		end, _ :=   strconv.Atoi(parts[1])

		availableIngredients[i] = Range{
			start: start,
			end:   end,
		}
	}
	sort.Slice(availableIngredients, func(i, j int) bool {
		return availableIngredients[i].start < availableIngredients[j].start
	})

	if part2 {
		return partTwo(availableIngredients)
	}
	return partOne(availableIngredients, lines[1])
}

func partOne(availableIngredients []Range, freshIngredientsStr string) int {
	result := 0
	freshIngredientsArray := strings.Split(freshIngredientsStr, "\n")
	freshIngredients := make([]int, len(freshIngredientsArray))
	for i, str := range freshIngredientsArray {
		freshIngredients[i], _ = strconv.Atoi(str)
	}
	sort.Ints(freshIngredients)

	startingSearchIndex := 0
    for _, freshIngredient := range freshIngredients {
		for i := startingSearchIndex; i < len(availableIngredients); i++ {
			availableIngredient := availableIngredients[i]
				
			if availableIngredient.start <= freshIngredient && availableIngredient.end >= freshIngredient {
				result++
				break
			}
			if availableIngredient.start > freshIngredient {
				break
			}
			startingSearchIndex = i
		}
	}
	return result
}

func partTwo(availableIngredients []Range) int {
	result := 0
	
	// Merge overlapping ranges
	merged := []Range{}
	
	for _, current := range availableIngredients {
		if len(merged) == 0 {
			merged = append(merged, current)
			continue
		}
		
		last := &merged[len(merged)-1]
		
		// If current range overlaps or is adjacent to the last merged range
		if current.start <= last.end+1 {
			// Extend the last range if needed
			if current.end > last.end {
				last.end = current.end
			}
		} else {
			// No overlap, add as new range
			merged = append(merged, current)
		}
	}
	
	// Calculate total coverage
	for _, r := range merged {
		result += r.end - r.start + 1
	}
	
	return result
}

func splitInput(input string) []string {
	lines := strings.Split(input, "\n\n")
	return lines
}