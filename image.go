package decker

import (
	"image"

	"github.com/corona10/goimagehash"
)

// Image is a wrapper around image.Image
// but with the path to the file added aswell
type Image struct {
	image.Image
	Path   string
	ID     uint64
	Hash   *goimagehash.ImageHash
	IsBest bool
}
