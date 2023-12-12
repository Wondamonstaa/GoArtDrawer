package main

/*
Name: Kiryl Baravikou
UIN: 656339218
Project 3: Go Interfaces
Date: 11/30/2023
Professor: Jon Solworth
Description:
The goal of the following project is to create a program that will draw the specified shapes on the screen.
The colors are specified in the colors' matrix. The program will be able to draw the following shapes:
	1. Rectangle
	2. Circle
	3. Triangle
The program will be able to draw the shapes on the screen and save the output as a .ppm file.
*/

import (
	"errors"
	"fmt"
	"math"
	"os"
)

// All possible colors needed for the assignment
var RED, GREEN, BLUE, YELLOW, ORANGE, PURPLE, BROWN, BLACK, WHITE Color

// Provided in the assignment description
var colors = [][]int{
	{255, 0, 0},
	{0, 255, 0},
	{0, 0, 255},
	{255, 255, 0},
	{255, 164, 0},
	{128, 0, 128},
	{165, 42, 42},
	{0, 0, 0},
	{255, 255, 255},
}

// Helper function to initialize the colors with indices to be able to access their values within the colors => DONE
func init() {
	RED = 0
	GREEN = 1
	BLUE = 2
	YELLOW = 3
	ORANGE = 4
	PURPLE = 5
	BROWN = 6
	BLACK = 7
	WHITE = 8
}

// Values used to represent the possible errors
var outOfBoundsErr error
var colorUnknownErr error

type Color int

// Point Object to represent a point
type Point struct {
	x, y int
}

// Triangle Object to represent the triangle
type Triangle struct {
	pt0, pt1, pt2 Point // three points to build the triangle
	c             Color
}

// Rectangle Object to represent the rectangle
type Rectangle struct {
	ll, ur Point // lower-left, upper-right
	c      Color
}

// Circle Object to represent the circle
type Circle struct {
	cp    Point // center point
	r     int   // radius
	color Color
}

// Display Object to represent the display
type Display struct {
	maxX, maxY   int
	screenColors [][]Color
}

// Screen Interface as specified in the assignment
type screen interface {
	getPixel(x, y int) (Color, error)
	background()
	screenShot(f string) error
	getMaxXY() (maxX, maxY int)
	initialize(maxX, maxY int)
	drawPixel(x, y int, c Color) error
}

// Geometry Interface as specified in the assignment
type geometry interface {
	draw(scn screen) error
	shape() string
}

// https://go.dev/blog/intro-generics
type Ordered interface {
	int | float64 | ~string
}

// Shape is a generic interface for shapes.
type Shape[T any] interface {
	Info() T
}

// FIXME: PROFESSOR
// https://gabrielgambetta.com/computer-graphics-from-scratch/07-filled-triangles.html
func interpolate(l0, d0, l1, d1 int) (values []int) {
	a := float64(d1-d0) / float64(l1-l0)
	d := float64(d0)

	count := l1 - l0 + 1
	for ; count > 0; count-- {
		values = append(values, int(d))
		d = d + a
	}
	return
}

// FIXME: PROFESSOR
// https://gabrielgambetta.com/computer-graphics-from-scratch/07-filled-triangles.html
func (tri Triangle) draw(scn screen) (err error) {
	if outOfBounds(tri.pt0, scn) || outOfBounds(tri.pt1, scn) || outOfBounds(tri.pt2, scn) {
		return outOfBoundsErr
	}
	if colorUnknown(tri.c) {
		return colorUnknownErr
	}

	y0 := tri.pt0.y
	y1 := tri.pt1.y
	y2 := tri.pt2.y

	// Sort the points so that y0 <= y1 <= y2
	if y1 < y0 {
		tri.pt1, tri.pt0 = tri.pt0, tri.pt1
	}
	if y2 < y0 {
		tri.pt2, tri.pt0 = tri.pt0, tri.pt2
	}
	if y2 < y1 {
		tri.pt2, tri.pt1 = tri.pt1, tri.pt2
	}

	x0, y0, x1, y1, x2, y2 := tri.pt0.x, tri.pt0.y, tri.pt1.x, tri.pt1.y, tri.pt2.x, tri.pt2.y

	x01 := interpolate(y0, x0, y1, x1)
	x12 := interpolate(y1, x1, y2, x2)
	x02 := interpolate(y0, x0, y2, x2)

	// Concatenate the short sides

	x012 := append(x01[:len(x01)-1], x12...)

	// Determine which is left and which is right
	var x_left, x_right []int
	m := len(x012) / 2
	if x02[m] < x012[m] {
		x_left = x02
		x_right = x012
	} else {
		x_left = x012
		x_right = x02
	}

	// Draw the horizontal segments
	for y := y0; y <= y2; y++ {
		for x := x_left[y-y0]; x <= x_right[y-y0]; x++ {
			err := scn.drawPixel(x, y, tri.c)
			if err != nil {
				return err
			}
		}
	}
	return
}

/* Display section */
var display Display

// Helper function to initialize the two points of the display => DONE
func (d *Display) initialize(maxX, maxY int) {
	d.maxX, d.maxY = maxX, maxY
	d.background()
}

// Professor's function to get the max X and Y dimensions of the screen => DONE
func (d *Display) getMaxXY() (maxX, maxY int) {
	return d.maxX, d.maxY
}

// Helper function to draw the pixel with color c at location x,y => DONE
func (d *Display) drawPixel(x, y int, c Color) error {
	return func() error {
		// Sanity check
		if outOfBounds(Point{x, y}, d) {
			return errors.New("Error: pixels out of bounds!")
		}
		// Initializes the selected color on the given x-y coordinates of the screen
		d.screenColors[x][y] = c
		return nil
	}()
}

// Helper function to get the color of the pixel at location x,y => DONE
func (d *Display) getPixel(x, y int) (Color, error) {
	return func() (Color, error) {
		//Bounds checking. The current size of the x, y must not exceed its max value
		if y < 0 || y >= d.maxY || x < 0 || x >= d.maxX || y < 0 || y >= d.maxY {
			return 0, errors.New("Error: pixels out of bounds!")
		}
		return d.screenColors[x][y], nil
	}()
}

/*
The following function sets the initial background colors for the Display => DONE
1. d => *Display - the display for which the background is set.
2. screenColors => [][]Color - the 2D slice to store pixel colors
*/
func (d *Display) background() {

	// Make allows to allocate and initialize the matrix of colors
	d.screenColors = make([][]Color, d.maxY)

	// Current row that will be filled with the specified color
	row := make([]Color, d.maxX)
	for i := range row {
		//Background => MUST be white
		row[i] = WHITE
	}

	// Finally, we fill the screenColors slice with the colored row
	for i := range d.screenColors {
		d.screenColors[i] = append(d.screenColors[i], row...)
	}
}

/*
Helper function from the project 2 to dump the screen as a .ppm file => DONE
1. d => *Display - the display for which the background is set.
2. f => string - the name of the file to write to.
*/
func (d *Display) screenShot(f string) error {

	//Create or truncate the named file in the .ppm format => also check for errors
	file, err := os.Create(f + ".ppm")
	if err != nil {
		return err
	}

	//Close the file once writing is done
	defer func(file *os.File) {
		err := file.Close()
		//Sanity check for successfully closing the file
		if err != nil {
			fmt.Println("Error: unable to close the given file!")
		}
	}(file)

	// Header creation that will be used within the generated .ppm file
	header := fmt.Sprintf("P3\n%d %d\n255\n", d.maxX, d.maxY)
	_, err = file.WriteString(header)

	//Sanity check again
	if err != nil {
		return err
	}

	// Finally, we need to populate the screen with the specified colors row by row
	for _, row := range d.screenColors {
		populateRow(file, row)
	}

	//Base case
	return nil
}

// Helper to return the info about the rectangle => DONE
func (r Rectangle) Info() string {
	return "Rectangle"
}

// Helper function to check if the rectangle's parameters are within bounds and have valid color => DONE
// @param s: screen => the screen on which the rectangle is drawn.
func (r Rectangle) sanityCheck(s screen) error {

	if func() bool {
		//Points and corners must be within the bound too to get the picture drawn
		checkBounds := func(p Point) bool {
			return outOfBounds(p, s)
		}
		return checkBounds(r.ll) || checkBounds(r.ur)
	}() {
		return outOfBoundsErr
	}
	// Valid range for the selected color
	if colorUnknown(r.c) {
		return colorUnknownErr
	}

	//Base case
	return nil
}

/*
Helper function to draw a filled rectangle on the current screen => DONE
From the Piazza:
1. ll => lower-left corner
2. ur => upper-right corner
3. c => fill color
@param s: screen - The screen on which the rectangle is drawn.
Reference: https://go.dev/doc/progs/image_draw.go
*/
func (r Rectangle) draw(s screen) error {

	//Sanity check
	if err := r.sanityCheck(s); err != nil {
		return err
	}

	//Helper function to draw a pixel on the screen with the specified color
	addPixels := func(x, y int) {
		err := s.drawPixel(x, y, r.c)
		//Must not be null
		if err != nil {
			return
		}
	}

	//Helper function to draw a row of pixels within the specified range
	addRows := func(y, startX, endX int) {
		for i := startX; i <= endX; i++ {
			addPixels(i, y)
		}
	}

	//Helper function to draw a column of pixels within the specified range
	addCols := func(startY, endY int) {
		for y := startY; y <= endY; y++ {
			addRows(y, r.ll.x, r.ur.x)
		}
	}

	//Extract rectangle parameters
	y0 := r.ll.y
	y1 := r.ur.y
	if y1 < y0 {
		y0, y1 = y1, y0
	}

	//Finally, draw the rectangle from the bottom to the top
	addCols(y0, y1)

	//Base case
	return nil
}

/*
Helper function to draw a filled circle on the specified screen => DONE
Info:
1. cp => center point
2. r => radius
3. color => fill color
@param s: screen => our screen on which the circle will be drawn
Reference: https://go.dev/doc/progs/image_draw.go
*/
func (c Circle) draw(s screen) error {

	// Sanity check for the correct bounds
	if outOfBounds(c.cp, s) {
		return outOfBoundsErr
	}

	// Color validation
	if colorUnknown(c.color) {
		return colorUnknownErr
	}

	//Unpacking the current parameters for the circle
	cX, cY, r2 := c.cp.x, c.cp.y, c.r
	maxX, _ := s.getMaxXY()

	//Helper function to draw a band of pixels within the specified range
	var band func(y, start, end int)
	band = func(y, start, end int) {
		if start <= end {
			err := s.drawPixel(start, y, c.color)
			if err != nil {
				return
			}
			band(y, start+1, end)
		}
	}

	// Recursively draws the circle from top to bottom
	var drawCircle func(y int)
	drawCircle = func(y int) {
		if y <= cY+r2 {
			bandWidth := int(math.Sqrt(math.Pow(float64(r2), 2) - math.Pow(float64(y-cY), 2)))
			startX := maxVal(cX-bandWidth, 0)
			endX := minVal(cX+bandWidth, maxX-1)
			band(y, startX, endX)
			drawCircle(y + 1)
		}
	}

	// Start drawing from the top of the circle
	drawCircle(cY - r2)

	//Base case
	return nil
}

// Helper to return the info about the circle => DONE
func (c Circle) shape() string {
	return "Circle"
}

// Helper function to find the min of two ints => DONE
func minVal(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Helper function to find the max of two ints => DONE
func maxVal(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Helper function to compute the pythagorean distance between two points => DONE
func pythagoreanDistance(x1, y1, x2, y2 int) float64 {
	dx := float64(x2 - x1)
	dy := float64(y2 - y1)

	return math.Pow(dx, 2) + math.Pow(dy, 2)
}

// Helper function used within the screenShot() to color the row with the specified colors => DONE
func populateRow(userFile *os.File, row []Color) {

	//Iterates over the row and colorizes the pixels with the specified colors from the colors matrix
	for _, c := range row {
		colorize(userFile, colors, c)
	}

	//Jump to the next line
	_, err := userFile.WriteString("\n")

	//Sanity check
	if err != nil {
		return
	}
}

// Helper function used within the populateRow() to colorize the pixels with the specified colors => DONE
func colorize(userFile *os.File, colors [][]int, c Color) {

	//Access the ppm and populate it with the specified colors at the given coordinates
	_, err := userFile.WriteString(fmt.Sprintf("%d %d %d ", colors[c][0], colors[c][1], colors[c][2]))

	//Sanity check
	if err != nil {
		return
	}
}

// Professor's function to check whether the point is out of bounds => DONE
func outOfBounds(p Point, s screen) bool {

	//Get the max values for both x and y
	maxX, maxY := s.getMaxXY()

	//Output message if error occurs => actually not the error => MUST be present there according to the pdf
	outOfBoundsErr = errors.New("geometry out of bounds")
	return p.y < 0 || p.y >= maxY || p.x < 0 || p.x >= maxX
}

// Professor's function to check whether the color is present in the provided pallet => DONE
func colorUnknown(c Color) bool {
	//Must be within the specified matrix of colors
	colorUnknownErr = errors.New("color unknown")
	var length = len(colors)
	return int(c) < 0 || int(c) >= length
}
