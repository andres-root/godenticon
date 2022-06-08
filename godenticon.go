package main

import (
	"crypto/md5"
	"image"
	"image/color"
	"log"

	"github.com/llgcode/draw2d/draw2dimg"
)

type GridPoint struct {
	value byte
	index int
}

type Identicon struct {
	name       string
	hash       [16]byte
	color      [3]byte
	grid       []byte
	gridPoints []GridPoint
	pixelMap   []DrawingPoint
}

type Point struct {
	x, y int
}

type DrawingPoint struct {
	topLeft     Point
	bottomRight Point
}

func getColor(hashFragment []byte) [3]byte {
	rgb := [3]byte{}
	copy(rgb[:], hashFragment[:])
	return rgb
}

func createGrid(hash [16]byte) []byte {
	grid := []byte{}

	for i := 0; i < len(hash) && i+3 <= len(hash)-1; i += 3 {
		chunk := make([]byte, 5)
		copy(chunk, hash[i:i+3])
		chunk[3] = chunk[1]
		chunk[4] = chunk[1]
		grid = append(grid, chunk...)

	}

	return grid
}

func filterOddSquares(grid []byte) []GridPoint {
	gridPoints := []GridPoint{}

	for i, code := range grid {
		if code%2 == 0 {
			point := GridPoint{
				value: code,
				index: i,
			}

			gridPoints = append(gridPoints, point)
		}
	}

	return gridPoints
}

func buildPixelMap(gridPoints []GridPoint) []DrawingPoint {
	drawingPoints := []DrawingPoint{}

	pixelFunc := func(p GridPoint) DrawingPoint {
		horizontal := (p.index % 5) * 50
		vertical := (p.index / 5) * 50
		topLeft := Point{horizontal, vertical}
		bottomRight := Point{horizontal + 50, vertical + 50}

		return DrawingPoint{
			topLeft,
			bottomRight,
		}
	}

	for _, gridPoint := range gridPoints {
		drawingPoints = append(drawingPoints, pixelFunc(gridPoint))
	}

	return drawingPoints
}

func rect(img *image.RGBA, col color.Color, x1, y1, x2, y2 float64) {
	gc := draw2dimg.NewGraphicContext(img) // Prepare new image context
	gc.SetFillColor(col)                   // set the color
	gc.MoveTo(x1, y1)                      // move to the topleft in the image
	// Draw the lines for the dimensions
	gc.LineTo(x1, y1)
	gc.LineTo(x1, y2)
	gc.MoveTo(x2, y1) // move to the right in the image
	// Draw the lines for the dimensions
	gc.LineTo(x2, y1)
	gc.LineTo(x2, y2)
	// Set the linewidth to zero
	gc.SetLineWidth(0)
	// Fill the stroke so the rectangle will be filled
	gc.FillStroke()
}

func drawRectangle(identiconColor [3]byte, pixelMap []DrawingPoint, name string) error {
	// We create our default image containing a 250x250 rectangle
	var img = image.NewRGBA(image.Rect(0, 0, 250, 250))
	// We retrieve the color from the color property on the identicon
	col := color.RGBA{identiconColor[0], identiconColor[1], identiconColor[2], 255}

	// Loop over the pixelmap and call the rect function with the img, color and the dimensions
	for _, pixel := range pixelMap {
		rect(
			img,
			col,
			float64(pixel.topLeft.x),
			float64(pixel.topLeft.y),
			float64(pixel.bottomRight.x),
			float64(pixel.bottomRight.y),
		)
	}
	// Finally save the image to disk
	return draw2dimg.SaveToPngFile(name+".png", img)
}

func createIdenticon(input []byte) Identicon {
	// Generate checksum from input
	name := string(input)
	hash := md5.Sum([]byte(input))
	color := getColor(hash[:3])
	grid := createGrid(hash)
	gridPoints := filterOddSquares(grid)
	pixelMap := buildPixelMap(gridPoints)

	return Identicon{
		name,
		hash,
		color,
		grid,
		gridPoints,
		pixelMap,
	}
}

func main() {
	data := []byte("bart")
	identicon := createIdenticon(data)

	if err := drawRectangle(identicon.color, identicon.pixelMap, identicon.name); err != nil {
		log.Fatalln(err)
	}
}
