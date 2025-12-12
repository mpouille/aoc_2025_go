package main

import (
	"github.com/jpillora/puzzler/harness/aoc"
	"strings"
	"strconv"
	"math"
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

	if (lines[0] == "..............") {
		fmt.Println("skip")
		return 0
	}

	coords := make([][]int, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		x, _ := strconv.Atoi(parts[0])
		y, _ := strconv.Atoi(parts[1])
		coords = append(coords, []int{x, y})
	}
	// when you're ready to do part 2, remove this "not implemented" block
	if part2 {
		return part2Run(coords)
	}
	// solve part 1 here
	return part1Run(coords)
}

func part1Run(coords [][]int) int {
	largestArea := 0

	for i := 0; i < len(coords); i++ {
		for j := i + 1; j < len(coords); j++ {
			area := calculateArea(coords[i], coords[j])
			if area > largestArea {
				largestArea = area
			}
		}
	}
	return largestArea
}

func calculateArea(coord1 []int, coord2 []int) int {
	width := int(math.Abs(float64(coord1[0] - coord2[0]))) + 1
	height := int(math.Abs(float64(coord1[1] - coord2[1]))) + 1
	return width * height
}

func part2Run(coords [][]int) int {
	// Convert coords to Point array
	points := make([]Point, len(coords))
	for i, coord := range coords {
		points[i] = Point{X: float64(coord[0]), Y: float64(coord[1])}
	}
	largestArea := 0 

	fmt.Println(len(points))
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			width := int(math.Abs(float64(points[i].X - points[j].X))) + 1
			height := int(math.Abs(float64(points[i].Y - points[j].Y))) + 1
			
			area :=  width * height
			fmt.Println(points[i], points[j], width, height, area, isRectangleInsidePolygon(points[i],  points[j], points))
			if area > largestArea && isRectangleInsidePolygon(points[i],  points[j], points) {
				fmt.Println("found")
				largestArea = area
			}
		}
	}	
	return largestArea
}

type Point struct {
	X, Y float64
}

// isInsidePolygon uses a robust ray casting algorithm (even-odd rule) that correctly 
// handles cases where the ray passes through vertices or along horizontal edges.
func isInsidePolygon(p Point, polygon []Point) bool {
	n := len(polygon)
	if n < 3 {
		return false
	}
	inside := false
	for i := 0; i < n; i++ {
		p1 := polygon[i]
		p2 := polygon[(i+1)%n]

		// Check if point is exactly a vertex
		if math.Abs(p.X-p1.X) < 1e-9 && math.Abs(p.Y-p1.Y) < 1e-9 {
			return true // Point is on the boundary/vertex, so it's inside
		}

		// Robust check for horizontal ray intersection
		if (p1.Y > p.Y) != (p2.Y > p.Y) {
			// Calculate the x-coordinate where the segment intersects the ray
			intersectX := (p2.X - p1.X) * (p.Y - p1.Y) / (p2.Y - p1.Y) + p1.X
			
			// If intersection is to the right of the point, toggle the inside status
			if p.X < intersectX {
				inside = !inside
			}
		}
	}
	return inside
}

// segmentsIntersect checks if two line segments strictly intersect (cross over),
// ignoring cases where they are collinear and overlapping, or just touching at endpoints.
// This function ensures that the rectangle edges do not cross *into* the polygon's interior
// from an exterior "cutout" area.
func segmentsIntersect(p1, q1, p2, q2 Point) bool {
	orientation := func(p, q, r Point) int {
		val := (q.Y-p.Y)*(r.X-q.X) - (q.X-p.X)*(r.Y-q.Y)
		if math.Abs(val) < 1e-9 { return 0 } // Collinear
		if val > 0 { return 1 } // Clockwise
		return 2 // Counterclockwise
	}

	o1 := orientation(p1, q1, p2)
	o2 := orientation(p1, q1, q2)
	o3 := orientation(p2, q2, p1)
	o4 := orientation(p2, q2, q1)

	// General case
	if o1 != o2 && o3 != o4 {
		return true // Strict crossing
	}
	// Note: Collinear cases that overlap are intentionally excluded here, 
	// as the user wants boundaries to be allowed to coincide.
	return false
}

// isRectangleInsidePolygon checks if an axis-aligned rectangle formed by two diagonal points (p1, p3)
// is entirely inside the given polygon.
func isRectangleInsidePolygon(p1, p3 Point, polygon []Point) bool {
	// Define the other two points of an axis-aligned rectangle
	p2 := Point{X: p1.X, Y: p3.Y}
	p4 := Point{X: p3.X, Y: p1.Y}

	rectanglePoints := []Point{p1, p2, p3, p4}
	
	// Create rectangle edges 
	rectEdges := [][2]Point{
		{p1, p2}, 
		{p2, p3}, 
		{p3, p4}, 
		{p4, p1},
	}

	// 1. Check if all four corners are inside/on the boundary of the polygon
	for _, p := range rectanglePoints {
		if !isInsidePolygon(p, polygon) {
			return false 
		}
	}

	// 2. Check if any rectangle edge strictly crosses any polygon edge
	n := len(polygon)
	for i := 0; i < n; i++ {
		polyP1 := polygon[i]
		polyP2 := polygon[(i+1)%n]

		for _, rectEdge := range rectEdges {
			if segmentsIntersect(rectEdge[0], rectEdge[1], polyP1, polyP2) {
				return false // A strict crossing intersection was found
			}
		}
	}

	// If all points are inside/on the boundary and no edges strictly cross the boundary, 
	// the entire rectangle is contained.
	return true
}