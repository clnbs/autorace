package models

import (
	"image"

	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
)

// Checkpoint are waypoint in the racetrack.
// -- currently not using it
type Checkpoint struct {
	CpSprite *pixel.Sprite
	Image    *pixel.PictureData
	Position pixel.Vec
	Number   int
}

// NewCheckpoint generate a graphical representation of a checkpoint
func NewCheckpoint(position pixel.Vec, number int) *Checkpoint {
	cp := new(Checkpoint)
	cp.Number = number
	cp.Position = position

	//Creating checkpoint image
	x0Size, y0Size, x1Size, y1Size := 0, 0, 20, 100
	img := image.NewRGBA(image.Rect(x0Size, y0Size, x1Size, y1Size))
	for indexY := y0Size; indexY < y1Size; indexY++ {
		for indexX := x0Size; indexX < x1Size; indexX++ {
			if (indexY/10)%2 == 0 && (indexX/10)%2 == 0 {
				img.Set(indexX, indexY, colornames.Black)
				continue
			}
			if (indexY/10)%2 == 1 && (indexX/10)%2 == 1 {
				img.Set(indexX, indexY, colornames.Black)
				continue
			}
			img.Set(indexX, indexY, colornames.White)
		}
	}
	cp.Image = pixel.PictureDataFromImage(img)
	cp.CpSprite = pixel.NewSprite(cp.Image, cp.Image.Bounds())
	return cp
}
