package models

import (
	"image"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"

	"golang.org/x/image/colornames"
)

// Car hold graphical players' car
type Car struct {
	CarSprite          *pixel.Sprite
	CarColor           color.RGBA
	Image              *pixel.PictureData
	CarName            *text.Text
	Angle              float64
	Position           pixel.Vec
}

// NewCar create a new car with a random color
func NewCar(carColor color.RGBA) *Car {
	newCar := new(Car)
	newCar.CarColor = carColor

	//Creating image for sprite
	x0Size, y0Size, x1Size, y1Size := 0, 0, 30, 20
	img := image.NewRGBA(image.Rect(x0Size, y0Size, x1Size, y1Size))
	for indexY := y0Size; indexY < y1Size; indexY++ {
		for indexX := x0Size; indexX < x1Size; indexX++ {
			if indexX > 20 {
				img.Set(indexX, indexY, colornames.Black)
				continue
			}
			img.Set(indexX, indexY, newCar.CarColor)
		}
	}
	//Creating sprite from image
	newCar.Image = pixel.PictureDataFromImage(img)
	newCar.CarSprite = pixel.NewSprite(newCar.Image, newCar.Image.Bounds())

	newCar.Angle = 0.0
	newCar.Position.Y = 0.0
	newCar.Position.X = 0.0

	return newCar
}
