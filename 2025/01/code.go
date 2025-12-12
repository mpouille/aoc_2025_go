package main

import (
	"strings"
	"fmt"
	"github.com/jpillora/puzzler/harness/aoc"
	"strconv"
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
	delta := map[byte]int{'L': -1, 'R': 1}
	dial := 50
	zeroCount := 0
	zeroCount2 := 0

	for i, s := range strings.Fields(string(input)) {
		n, _ := strconv.Atoi(s[1:])

		for range n {
			if dial += delta[s[0]]; dial%100 == 0 {
				fmt.Printf("ZERO FOUND line %d: %s\n", i, s)
				zeroCount2++
			}
		}
		if dial%100 == 0 {
			zeroCount++
		}
	}

	if part2 {
		return zeroCount2
	}

	// solve part 1 here
	return zeroCount
}


