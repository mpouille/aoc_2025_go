package main

import (
	"github.com/jpillora/puzzler/harness/aoc"
	"github.com/draffensperger/golp"
	"math"
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
	machines := parseInput(input)
	
	// when you're ready to do part 2, remove this "not implemented" block
	if part2 {
		return part2Run(machines)
	}
	// solve part 1 here
	return part1Run(machines)
}

func part1Run(machines []Machine) int {
	result := 0
	for _, machine := range machines {
		result += calculateCostForMachine(machine)
	}
	return result
}

func calculateCostForMachine(machine Machine) int {
	lights := machine.lightDiagram
	buttons := machine.buttons
	seen := make(map[int]bool)
	seen[0] = true
	
	costs := [][]int{{0}}
	
	for curCost := 1; ; curCost++ {
		oldSeenLength := len(seen)
		
		prevCosts := costs[len(costs)-1]
		
		var curCosts []int
		
		for _, button := range buttons {
			for _, parent := range prevCosts {
				candidate := button ^ parent
				
				if seen[candidate] {
					continue
				}
				
				if candidate == lights {
					return curCost
				}
				
				seen[candidate] = true
				curCosts = append(curCosts, candidate)
			}
		}
		
		costs = append(costs, curCosts)
		
		if len(seen) <= oldSeenLength {
			panic("assertion failed: seen length should increase")
		}
	}
}

func part2Run(machines []Machine) int {
	total := 0
	for _, machine := range machines {
		total += minimizeInputs(machine.joltageRequirements, machine.buttonsPart2)
	}
	return total
}

type Machine struct {
	lightDiagram int
	buttons []int
	buttonsPart2 [][]int
	joltageRequirements []int
}

func parseInput(input string) []Machine {
	lines := strings.Split(input, "\n")
	
	machines := make([]Machine, 0, len(lines))

	for _, line := range lines {
		machine := Machine{}
	
		sections := strings.Split(line, " ")

		// Process light diagram [#.#.] at index 0
		runes := []rune(strings.NewReplacer(".", "0", "#", "1").Replace(strings.Trim(sections[0], "[]")))
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		val, _ := strconv.ParseInt(string(runes), 2, 64)
		machine.lightDiagram = int(val)
		// Process buttons (1,2) (3,4,5) from index 1 to len(sections) - 2
		machine.buttons = make([]int, 0, len(sections) - 2)
		machine.buttonsPart2 = make([][]int, 0, len(sections) - 2)
		for i := 1; i < len(sections) - 1; i++ {
			buttonStrArray := strings.Split(strings.Trim(sections[i], "()"), ",")
			buttonMask := 0
			buttonMaskPart2 := make([]int, 0, len(buttonStrArray))
			for _, buttonStr := range buttonStrArray {
				// Convert the index string to an integer
				index, _ := strconv.Atoi(buttonStr)
				// Create the bitmask: 1 << index is 2^index
				// This sets the bit at the position indicated by the index.
				buttonMask |= (1 << index)
				buttonMaskPart2 = append(buttonMaskPart2, index)
			}
			machine.buttons = append(machine.buttons, buttonMask)
			machine.buttonsPart2 = append(machine.buttonsPart2, buttonMaskPart2)
		}
		// Process joltage requirements {1,2,3} from index len(sections) - 1
		joltageRequirementsStrArray := strings.Split(strings.Trim(sections[len(sections) - 1], "{}"), ",")
		machine.joltageRequirements = make([]int, 0, len(joltageRequirementsStrArray))
		for _, joltageRequirementStr := range joltageRequirementsStrArray {
			joltageRequirementInt, _ := strconv.Atoi(joltageRequirementStr)
			machine.joltageRequirements = append(machine.joltageRequirements, joltageRequirementInt)
		}
		machines = append(machines, machine)
	}

	return machines
}

func minimizeInputs(joltages []int, buttonIndexes [][]int) int {
	numButtons := len(buttonIndexes)
	numJoltages := len(joltages)
	
	// Build LP: variables = presses per button (integer, >= 0)
	lp := golp.NewLP(0, numButtons)
	// Objective: minimize sum(x_i)
	obj := make([]float64, numButtons)
	for i := 0; i < numButtons; i++ {
		obj[i] = 1.0
	}
	lp.SetObjFn(obj)
	// Variable types and bounds
	for i := 0; i < numButtons; i++ {
		lp.SetInt(i, true)              // integer variables
		lp.SetBounds(i, 0.0, math.Inf(1)) // x_i >= 0
	}
	// Constraints: for each joltage j: sum over affecting buttons x_i = joltages[j]
	for j := 0; j < numJoltages; j++ {
		entries := make([]golp.Entry, 0, numButtons)
		for i := 0; i < numButtons; i++ {
			affects := false
			for _, idx := range buttonIndexes[i] {
				if idx == j {
					affects = true
					break
				}
			}
			if affects {
				entries = append(entries, golp.Entry{Col: i, Val: 1.0})
			}
		}
		target := float64(joltages[j])
		if len(entries) == 0 {
			if target != 0 {
				return -1 // unsatisfiable: no buttons affect this counter
			}
			continue // skip constraint 0 = 0
		}
		_ = lp.AddConstraintSparse(entries, golp.EQ, target)
	}
	// Solve
	res := lp.Solve()
	if res != golp.OPTIMAL {
		return -1
	}
	return int(math.Round(lp.Objective()))
}
