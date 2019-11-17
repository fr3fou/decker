package decker

import "image"

// Image is a wrapper around image.Image
// but with the path to the file added aswell
type Image struct {
	image.Image
	Path     string
	Hash     uint64
	IsBest   bool
	Siblings []Image
}
