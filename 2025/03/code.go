package main

import (
	"github.com/jpillora/puzzler/harness/aoc"
	"strings"
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

	lines := strings.Split(strings.TrimSpace(input), "\n")	

	// when you're ready to do part 2, remove this "not implemented" block
	if part2 {
		return part2Fn(lines)
	}
	
	return part1(lines)
}

func part1(lines []string) int {
	result := 0
	for _, line := range lines {
		maxJoltage, _ := strconv.Atoi(recursiveFindHighestJoltage(line, 2))
		result += maxJoltage
	}
	return result
}

func part2Fn(lines []string) int {
	result := 0
	for _, line := range lines {
		// recursively find the highest joltage ( with length - 12)
		maxJoltage, _ := strconv.Atoi(recursiveFindHighestJoltage(line, 12))
		result += maxJoltage
	}
	return result
}

func recursiveFindHighestJoltage(line string, length int) string {
	if length == 0 {
		return ""
	}
	if length >= len(line)  {
		return line
	}

	maxJoltage := 0
	maxJoltageIndex := 0
	for i := 0; i < len(line) - length + 1; i++ {
		if line[i] >= '0' && line[i] <= '9' {
			digit := int(line[i] - '0')
			if digit > maxJoltage {
				maxJoltage = digit
				maxJoltageIndex = i
			}
		}
	}

	return strconv.Itoa(maxJoltage) + recursiveFindHighestJoltage(line[maxJoltageIndex + 1:], length-1)
}