package main

import (
	"github.com/jpillora/puzzler/harness/aoc"
	"strings"
	"strconv"	
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
	problemCount := len(strings.Fields(lines[0]))
	linesCount := len(lines)
	problems := make([][]string, linesCount)
	for i, line := range lines {
		problems[i] = make([]string, problemCount)
		parts := strings.Fields(line)
		for j, str := range parts {
			problems[i][j] = str
		}
	}
	// when you're ready to do part 2, remove this "not implemented" block
	if part2 {
		return runPart2(lines);
	}
	// solve part 1 here
	return runPart1(problems);
}

func runPart1(problems [][]string) int {
	result := 0
	
	for i := 0; i < len(problems[0]); i++ {
		sum := 0
		operation := problems[len(problems) -1][i]
		for j := 0; j < len(problems) - 1; j++ {
			value, _ := strconv.Atoi(problems[j][i])
			if operation == "+" {
				sum += value
			} else if operation == "*" {
				if j == 0 {
					sum = 1
				}
				sum *= value
			}
		}
		result += sum
	}
	return result
}

func runPart2(lines []string) int {
	// Split each line into individual characters
	charGrid := make([][]string, len(lines))
	maxLength := 0
	for i, line := range lines {
		charGrid[i] = make([]string, len(line))
		maxLength = max(maxLength, len(line))
		for j, char := range line {
			charGrid[i][j] = string(char)
		}
	}

    currentOperation := -1
	currentOperationResult := 0
	result := 0

	for i := 0; i < maxLength; i++ {
		currentNumber := ""
		for j := 0; j < len(charGrid); j++ {
			if i >= len(charGrid[j]) {
				continue
			}
			currentNumber += charGrid[j][i]
		}
        
		if strings.TrimSpace(currentNumber) == "" {
			continue
		}

		lastChar := currentNumber[len(currentNumber) - 1]
		if lastChar == '+' || lastChar == '*' {
			// new operation, save previous operation result
			fmt.Println(lastChar, " found, add ", currentOperationResult, " to result", result)
			result += currentOperationResult

			cephalNumber := currentNumber[:len(currentNumber)-1]
			num, _ := strconv.Atoi(strings.TrimSpace(cephalNumber))
			fmt.Println("Set current to ", cephalNumber)
			currentOperationResult = num
			if lastChar == '+' {
				currentOperation = 0
			} else {
				currentOperation = 1
			}
		} else {
			cephalNumber, _ := strconv.Atoi(strings.TrimSpace(currentNumber))
			fmt.Println("Set current to ", cephalNumber)
			if currentOperation == 0 {
				currentOperationResult += cephalNumber
			} else if currentOperation == 1 {
				currentOperationResult *= cephalNumber
			}
		}
		if i == maxLength - 1 {
			fmt.Println("Add ", currentOperationResult, " to result", result)
			result += currentOperationResult
		}
	}

	return result
}