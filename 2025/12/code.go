package main

import (
	"github.com/jpillora/puzzler/harness/aoc"
	"strings"
	"strconv"
	"fmt"
	"sort"
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
	presents, regions := parseInput(input)
	
	// when you're ready to do part 2, remove this "not implemented" block
	if part2 {
		return part2Run(presents, regions)
	}
	// solve part 1 here
	return part1Run(presents, regions)
}

func part1Run(presents []PresentShape, regions []Region) int {
	result := 0
	for _, region := range regions {
		// get presents from the region
		presentsInRegion := []PresentInRegion{}
		maxRequiredArea := 0
		for i, shapeCount := range region.shapesToFit {
			if shapeCount > 0 {
				presentsInRegion = append(presentsInRegion, PresentInRegion{presents[i], shapeCount})
				maxRequiredArea += presents[i].shapeSize * shapeCount
			}
		}

		if maxRequiredArea <= region.width * region.length {
			// init region grid
			regionGrid := make([][]bool, region.width)
			for i := range regionGrid {
				regionGrid[i] = make([]bool, region.length)
			}
			// Try to fit the shapes in the region.
			if SolvePacking(presentsInRegion, regionGrid).Success {
				result++
			}
		}
	}
	return result
}

func part2Run(presents []PresentShape, regions []Region) int {
	return 42
}

type Placement struct {
    shapeIndex string
    orientation [][]bool // The specific orientation being used
    x, y int             // The position where it was placed
}

// Result struct to communicate the outcome of the packing attempt.
type PackingResult struct {
    Success bool
    PlacedPlacements []Placement // List of all successfully placed items
    FinalRegion [][]bool
    UnplacedShapes []PlaceableShape // List of shapes that couldn't be fitted
}

// Helper struct for shapes during the sorting/placement phase
type PlaceableShape struct {
    // Note: We need a unique ID for tracking, so we'll use the original shape index + instance ID
    instanceID string // e.g., "ShapeA_1", "ShapeA_2"
    orientations [][][]bool
    size int
}

// Get all the distinct shape orientations from the PresentShape struct
func GetOrientations(ps PresentShape) [][][]bool {
    // Collect all valid orientations into a slice
    return [][][]bool{
        ps.shape, 
        ps.shapeRotated90, ps.shapeRotated180, ps.shapeRotated270,
        ps.shapeFlipped, 
        ps.shapeFlippedRotated90, ps.shapeFlippedRotated180, ps.shapeFlippedRotated270,
    }
    // Note: You might need to add logic to filter out duplicates if the shape
    // is symmetric (e.g., a square 2x2 has only 1 unique orientation).
}

type PresentInRegion struct {
	present PresentShape
	count int
}

type PresentShape struct {
	index int
	shape [][]bool
	shapeRotated90 [][]bool
	shapeRotated180 [][]bool
	shapeRotated270 [][]bool
	shapeFlipped [][]bool
	shapeFlippedRotated90 [][]bool
	shapeFlippedRotated180 [][]bool
	shapeFlippedRotated270 [][]bool
	shapeSize int
}

func rotateShape90(shape [][]bool) [][]bool {
	newShape := make([][]bool, len(shape[0]))
	for i := range newShape {
		newShape[i] = make([]bool, len(shape))
	}
	for i := range shape {
		for j := range shape[i] {
			newShape[j][len(shape) - i - 1] = shape[i][j]
		}
	}
	return newShape
}

func flipShape(shape [][]bool) [][]bool {
	newShape := make([][]bool, len(shape))
	for i := range newShape {
		newShape[i] = make([]bool, len(shape[i]))
	}
	for i := range shape {
		for j := range shape[i] {
			newShape[i][j] = shape[i][len(shape[i]) - j - 1]
		}
	}
	return newShape
}

type Region struct {
	width int
	length int
	shapesToFit []int
}

func parseInput(input string) ([]PresentShape, []Region) {
	sections := strings.Split(input, "\n\n")

	presents := make([]PresentShape, len(sections) - 1)
	for i, section := range sections[:len(sections) - 1] {
		lines := strings.Split(section, "\n")

		shape := make([][]bool, len(lines[1:]))
		shapeSize := 0
		for i, line := range lines[1:] {
			shape[i] = []bool{}
			for _, char := range line {
				shape[i] = append(shape[i], char == '#')
			}
			shapeSize += strings.Count(line, "#")
		}
		// generate all rotated and flipped shapes to ease checks later on
		shapeRotated90 := rotateShape90(shape)
		shapeRotated180 := rotateShape90(shapeRotated90)
		shapeRotated270 := rotateShape90(shapeRotated180)
		shapeFlipped := flipShape(shape)
		shapeFlippedRotated90 := rotateShape90(shapeFlipped)
		shapeFlippedRotated180 := rotateShape90(shapeFlippedRotated90)
		shapeFlippedRotated270 := rotateShape90(shapeFlippedRotated180)

		index, _ := strconv.Atoi(strings.Trim(lines[0], ":"))
		presents[i] = PresentShape{
			index: index,
			shape: shape,
			shapeSize: shapeSize,
			shapeRotated90: shapeRotated90,
			shapeRotated180: shapeRotated180,
			shapeRotated270: shapeRotated270,
			shapeFlipped: shapeFlipped,
			shapeFlippedRotated90: shapeFlippedRotated90,
			shapeFlippedRotated180: shapeFlippedRotated180,
			shapeFlippedRotated270: shapeFlippedRotated270,
		}
	}

	// regions: 1234x5678: 1 2 3 4
	regionStrs := strings.Split(sections[len(sections) - 1], "\n")
	regions := make([]Region, len(regionStrs))
	for i, regionStr := range regionStrs {
		lines := strings.Split(regionStr, "\n")
		for _, line := range lines {
			lineSplit := strings.Split(line, ": ")
			dims := strings.Split(lineSplit[0], "x")
			width, _ := strconv.Atoi(dims[0])
			length, _ := strconv.Atoi(dims[1])
			
			shapesToFitStrs := strings.Split(lineSplit[1], " ")
			shapesToFit := []int{}
			for _, shapeIndex := range shapesToFitStrs {
				shapeIndexInt, _ := strconv.Atoi(shapeIndex)
				shapesToFit = append(shapesToFit, shapeIndexInt)
			}

			regions[i] = Region{
				width: width,
				length: length,
				shapesToFit: shapesToFit,
			}
		}
	}

	return presents, regions
}

func printShape(shape [][]bool) {
	for _, row := range shape {
		for _, cell := range row {
			if cell {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func CanPlace(region [][]bool, shapeOrientation [][]bool, x int, y int) bool {
    shapeH := len(shapeOrientation)
    shapeW := len(shapeOrientation[0])
    regionH := len(region)
    regionW := len(region[0])

    // Iterate through the shape's local coordinates (i, j)
    for i := 0; i < shapeH; i++ {
        for j := 0; j < shapeW; j++ {
            // Check if the current shape pixel is 'true' (part of the shape)
            if shapeOrientation[i][j] {
                // Calculate the corresponding absolute coordinates (absX, absY)
                absX := x + j
                absY := y + i 

                // 1. Check Out-of-Bounds
                if absY < 0 || absY >= regionH || absX < 0 || absX >= regionW {
                    // This pixel is out of bounds
                    return false
                }

                // 2. Check Overlap (Collision)
                if region[absY][absX] {
                    // This pixel is already occupied in the region
                    return false
                }
            }
        }
    }
    // No collision or out-of-bounds detected
    return true
}

// SolvePacking attempts to place all shapes in the region using a greedy approach.
func SolvePacking(shapesToFit []PresentInRegion, region [][]bool) PackingResult {
    regionH := len(region)
    regionW := len(region[0])
    currentRegion := make([][]bool, regionH)
    
    // Deep copy the initial region (which should be empty 'false' values)
    for i := range region {
        currentRegion[i] = make([]bool, regionW)
        copy(currentRegion[i], region[i])
    }

    // --- 1. Preprocess and Sort Shapes ---
    var allShapes []PlaceableShape
    
    for _, pir := range shapesToFit {
        orientations := GetOrientations(pir.present)
        
        // Add the shape 'count' number of times with unique instance IDs
        for i := 0; i < pir.count; i++ {
            allShapes = append(allShapes, PlaceableShape{
                instanceID: fmt.Sprintf("%d_%d", pir.present.index, i), // Unique ID tracking
                orientations: orientations,
                size: pir.present.shapeSize,
            })
        }
    }

    // Sort the list: Largest shapes first (Greedy choice)
    sort.Slice(allShapes, func(i, j int) bool {
        return allShapes[i].size > allShapes[j].size
    })

    // --- 2. Greedy Placement Loop ---
    var placedPlacements []Placement
    var unplacedShapes []PlaceableShape

    for _, pShape := range allShapes {
        bestPlacementFound := false

        // Simple iteration over all coordinates (x, y)
        for y := 0; y < regionH; y++ {
            for x := 0; x < regionW; x++ {
                
                // Iterate through all possible orientations
                for _, orientation := range pShape.orientations {
                    
                    if CanPlace(currentRegion, orientation, x, y) {
                        // --- Commit Placement ---
                        shapeH := len(orientation)
                        shapeW := len(orientation[0])
                        
                        // Mark cells as occupied
                        for i := 0; i < shapeH; i++ {
                            for j := 0; j < shapeW; j++ {
                                if orientation[i][j] {
                                    currentRegion[y+i][x+j] = true
                                }
                            }
                        }

                        // Record the placement
                        placedPlacements = append(placedPlacements, Placement{
                            shapeIndex: pShape.instanceID, // Use the unique ID
                            orientation: orientation,
                            x: x, 
                            y: y,
                        })
                        
                        bestPlacementFound = true
                        goto NextShape // Jump out of all inner loops
                    }
                }
            }
        }

        NextShape: // Label to jump to when a placement is found
        if !bestPlacementFound {
            // This shape could not be fitted
            unplacedShapes = append(unplacedShapes, pShape)
        }
    }
    
    // --- 3. Return Final Result ---
    
    return PackingResult{
        Success: len(unplacedShapes) == 0,
        PlacedPlacements: placedPlacements,
        FinalRegion: currentRegion,
        UnplacedShapes: unplacedShapes,
    }
}