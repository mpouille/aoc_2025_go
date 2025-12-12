package main

import (
	"github.com/jpillora/puzzler/harness/aoc"	
	"math"
	"strings"
	"strconv"
	"slices"
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
	isExample := len(lines) == 20;

	points := make([][]int, len(lines))
	for i, line := range lines {
		parts := strings.Split(line, ",")
		x, _ := strconv.Atoi(parts[0])
		y, _ := strconv.Atoi(parts[1])
		z, _ := strconv.Atoi(parts[2])
		points[i] = []int{x, y, z}
	}

	// when you're ready to do part 2, remove this "not implemented" block
	if part2 {
		return part2Run(points)
	}
	// solve part 1 here
	return part1Run(points, isExample)
}

func part2Run(points [][]int) int {
	var edges []Edge

	// Generate all possible edges between points
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			p1 := points[i]
			p2 := points[j]
			distance := euclideanDistance(p1[0], p1[1], p1[2], p2[0], p2[1], p2[2])
			edges = append(edges, Edge{from: i, to: j, weight: distance})
		}
	}

	slices.SortFunc(edges, func(a, b Edge) int {
		return a.weight - b.weight
	})
	
	// Use binary search (dichotomy) to find minimum number of edges needed
	left := 1
	right := len(edges)

	for left <= right {
		mid := (left + right) / 2
		
		// Test with first 'mid' edges
		testEdges := edges[:mid]
		
		if checkIfGraphAllConnected(points, testEdges) {
			right = mid - 1  // Try with fewer edges
		} else {
			left = mid + 1   // Need more edges
		}
	}

	fmt.Println(points[edges[right].from],)

	return points[edges[right].from][0] * points[edges[right].to][0]
}

func checkIfGraphAllConnected(points [][]int, edges []Edge) bool {
	var graphs []Graph
	usedEdges := make(map[int]bool)

	for i, edge := range edges {
		if usedEdges[i] {
			continue
		}
		
		// Create a new graph starting with this edge
		graph := Graph{}
		graph.AddEdge(edge.from, edge.to, edge.weight)
		usedEdges[i] = true
		
		// Find all other edges that can link to this graph
		changed := true
		for changed {
			changed = false
			for j, otherEdge := range edges {
				if usedEdges[j] {
					continue
				}
				
				// Check if this edge can link to the current graph
				if graph.CanLinkToEdge(otherEdge.from, otherEdge.to) {
					graph.AddEdge(otherEdge.from, otherEdge.to, otherEdge.weight)
					usedEdges[j] = true
					changed = true
				}
			}
		}
		
		graphs = append(graphs, graph)
	}

	return len(graphs) == 1 && len(graphs[0].points) == len(points)
}

func part1Run(points [][]int, isExample bool) int {

	maxIterations := 1000;
	if isExample {
		maxIterations = 10;
	}
	var edges []Edge
	// Create a dynamic array to store up to ten edges
	smallestEdges := make([]Edge, 0, maxIterations)

	// Generate all possible edges between points
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			p1 := points[i]
			p2 := points[j]
			distance := euclideanDistance(p1[0], p1[1], p1[2], p2[0], p2[1], p2[2])
			edges = append(edges, Edge{from: i, to: j, weight: distance})
			if len(smallestEdges) < maxIterations {
				smallestEdges = append(smallestEdges, Edge{from: i, to: j, weight: distance})
			} else {
				// Find the edge with the greatest weight
				greatestWeight := 0
				greatestWeightIndex := 0
				for k, edge := range smallestEdges {
					if edge.weight > greatestWeight {
						greatestWeight = edge.weight
						greatestWeightIndex = k
					}
				}
				if smallestEdges[greatestWeightIndex].weight > distance {
					smallestEdges[greatestWeightIndex] = Edge{from: i, to: j, weight: distance}
				}
			}
		}
	}
	
	// Create graphs by linking edges that share endpoints
	var graphs []Graph
	usedEdges := make(map[int]bool)
	
	for i, edge := range smallestEdges {
		if usedEdges[i] {
			continue
		}
		
		// Create a new graph starting with this edge
		graph := Graph{}
		graph.AddEdge(edge.from, edge.to, edge.weight)
		usedEdges[i] = true
		
		// Find all other edges that can link to this graph
		changed := true
		for changed {
			changed = false
			for j, otherEdge := range smallestEdges {
				if usedEdges[j] {
					continue
				}
				
				// Check if this edge can link to the current graph
				if graph.CanLinkToEdge(otherEdge.from, otherEdge.to) {
					graph.AddEdge(otherEdge.from, otherEdge.to, otherEdge.weight)
					usedEdges[j] = true
					changed = true
				}
			}
		}
		
		graphs = append(graphs, graph)
	}
	
	// order graphs by size descending
	slices.SortFunc(graphs, func(a, b Graph) int {
		return len(b.points) - len(a.points)
	})

	// multiply the sizes of the three largest graphs
	return len(graphs[0].points) * len(graphs[1].points) * len(graphs[2].points)
}

func euclideanDistance(x1 int, y1 int, z1 int, x2 int, y2 int, z2 int) int {
	return int(math.Sqrt(float64((x1 - x2) * (x1 - x2) + (y1 - y2) * (y1 - y2) + (z1 - z2) * (z1 - z2))))
}

type Edge struct {
	from, to int
	weight int
}

type Graph struct {
	points []int
	edges []Edge
}

func (g *Graph) CanLinkToEdge(p1 int, p2 int) bool {
	return slices.Contains(g.points, p1) || slices.Contains(g.points, p2)
}

func (g *Graph) AddEdge(from, to int, weight int) {
	if !slices.Contains(g.points, from) {
		g.points = append(g.points, from)
	}
	if !slices.Contains(g.points, to) {
		g.points = append(g.points, to)
	}
	g.edges = append(g.edges, Edge{from: from, to: to, weight: weight})
}

func (g *Graph) GetEdges() []Edge {
	return g.edges
}
