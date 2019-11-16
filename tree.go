package decker

import "image"

// Tree is the struct that holds
// the best quality version of the image
// as well as any duplicates of it as children
type Tree struct {
	Best       Image
	Duplicates []Image
}

// Image is a wrapper around image.Image
// but with the path to the file added aswell
type Image struct {
	image.Image
	Path string
}
