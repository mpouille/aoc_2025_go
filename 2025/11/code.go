package main

import (
	"github.com/jpillora/puzzler/harness/aoc"	
	"strings"
	"fmt"
	"slices"
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

	// each line defines a node and its connections (format aaa: bbb ccc)
	graph := make(map[string][]string)
	for _, line := range lines {
		parts := strings.Split(line, ": ")
		graph[parts[0]] = strings.Split(parts[1], " ")
	}

	// when you're ready to do part 2, remove this "not implemented" block
	if part2 {
		return part2Run(graph)
	}
	// solve part 1 here
	return part1Run(graph)
}

func part1Run(graph map[string][]string) int {
	paths := pathsBetween(graph, make(map[string]int), "you", "out")
	return paths
}

func part2Run(graph map[string][]string) int {
	cachedPaths := make(map[string]int)
	// The original logic finds which midpoint ("dac" or "fft") is encountered first via BFS from "svr"
	queue := []string{"svr"}
	// Keep track of visited nodes during the BFS to prevent infinite loops if graph has cycles
	bfsVisited := make(map[string]bool) 
	bfsVisited["svr"] = true 

	firstMidpoint := ""

	for len(queue) > 0 {
		device := queue[0]
		queue = queue[1:] // Dequeue

		if device == "dac" || device == "fft" {
			firstMidpoint = device
			break // Exit the BFS loop once found
		}

		for _, output := range graph[device] {
			if !bfsVisited[output] {
				bfsVisited[output] = true
				queue = append(queue, output) // Enqueue neighbors
			}
		}
	}

    // Determine the second midpoint based on the first one found
	secondMidpoint := ""
	if firstMidpoint == "dac" {
		secondMidpoint = "fft"
	} else if firstMidpoint == "fft" {
		secondMidpoint = "dac"
	} else {
        // Handle case where neither midpoint was reachable (error handling)
        fmt.Println("Error: Neither 'dac' nor 'fft' found via BFS from 'svr'")
        return 0
    }

	part2Paths := pathsBetween(graph, cachedPaths, "svr", firstMidpoint) *
		pathsBetween(graph, cachedPaths, firstMidpoint, secondMidpoint) *
		pathsBetween(graph, cachedPaths, secondMidpoint, "out")

	return part2Paths
}

// Old way for part 1 through pure DFS
// While logging found path on part two, no loop were detected.
// Since pure DFS was way too slow, I went for memoized DFS (only possible in non cyclic graphs called "DAG")
func findAllPaths(graph map[string][]string, source, destination string) [][]string {
	var allPaths [][]string
	// Use a map to track visited nodes within the current path.
	// We use a map[string]bool instead of just a slice to make lookups faster (O(1)).
	visited := make(map[string]bool)
	// Start the recursive DFS from the source.
	dfs(graph, source, destination, []string{}, visited, &allPaths)
	return allPaths
}

// DFS to find a path from a to b
func dfs(graph map[string][]string, current, destination string, currentPath []string, visited map[string]bool, allPaths *[][]string) {
	if slices.Contains(currentPath, current) {
		fmt.Println("Loop detected: ", current)
		return
	}
	// 1. Mark the current node as visited.
	visited[current] = true
	// 2. Add the current node to the path.
	currentPath = append(currentPath, current)

	// 3. Base case: If the current node is the destination, a path is found.
	if current == destination {
		// Create a copy of the current path and append it to the results.
		// We copy because currentPath is a slice that gets modified in subsequent recursions.
		pathCopy := make([]string, len(currentPath))
		copy(pathCopy, currentPath)
		*allPaths = append(*allPaths, pathCopy)
	} else {
		// 4. Recurse for all unvisited neighbors.
		for _, neighbor := range graph[current] {
			if !visited[neighbor] {
				dfs(graph, neighbor, destination, currentPath, visited, allPaths)
			}
		}
	}

	// 5. Backtrack: Unmark the current node as visited and remove it from the path
	//    so that it can be included in other potential paths.
	currentPath = currentPath[:len(currentPath)-1]
	visited[current] = false
}

// As explained above, switch to an optimized version of DFS that counts paths from start to end using memoization
// Then to optimize further, the path calculation is split in 3 parts:
// svr -> firstMidpoint, 
// firstMidpoint -> secondMidpoint, 
// secondMidpoint -> out
// where firstMidpoint is the closest from svr (or the start essentially) between dac and fft
// and secondMidpoint is the other one
// Then the solution is the product of the 3 parts.
func pathsBetween(graph map[string][]string, cachedPaths map[string]int, start, end string) int {
	if start == end {
		return 1
	}

	if start == "out" { 
		return 0
	}

    // Create a unique cache key
	key := start + "-" + end 

    // Check memoization cache
	if paths, found := cachedPaths[key]; found {
		return paths
	}

	paths := 0
    // Recurse over neighbors
	for _, output := range graph[start] {
		paths += pathsBetween(graph, cachedPaths, output, end)
	}

    // Store the result in the cache before returning
	cachedPaths[key] = paths
	return paths
}

