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

	result := 0

	for _, s := range strings.Fields(string(input)) {
		// Split the string by commas to get individual ranges
		ranges := strings.Split(s, ",")
		
		for _, rangeStr := range ranges {
			// Split each range by hyphen to get start and end
			parts := strings.Split(rangeStr, "-")

			if len(parts) != 2 {
				continue
			}

			firstId, _ := strconv.Atoi(parts[0])
			lastId, _ := strconv.Atoi(parts[1])

			// Check each ID in the range for numbers displayed twice
			for id := firstId; id <= lastId; id++ {
				if id < 10 {
					continue
				}
				idStr := strconv.Itoa(id)
				
				// Part 1: 
				// Check if the ID is composed of a number displayed twice
				// This means the string length should be even and first half equals second half
				if !part2 && len(idStr)%2 == 0 {
					mid := len(idStr) / 2
					firstHalf := idStr[:mid]
					secondHalf := idStr[mid:]
					
					if firstHalf == secondHalf {
						result += id
					}
				}
				// Part 2:
				// Check if the ID is composed of a number displayed AT LEAST twice omg
				if part2 {
					// Check if the ID is composed of a number displayed at least twice
					// Find the smallest digit length that divides the string length evenly
					for d := 1; d < len(idStr); d++ {
						if len(idStr)%d != 0 {
							continue
						}
						pattern := idStr[:d]
						isRepeated := true

						for i := d; i < len(idStr); i += d {
							if i+d > len(idStr) || idStr[i:i+d] != pattern {
								isRepeated = false
								break
							}
						}
						if isRepeated {
							result += id
							break
						}
					}
				}
			}
		}
	}
	// when you're ready to do part 2, remove this "not implemented" block
	if part2 {
		return result
	}
	// solve part 1 here
	return result
}
