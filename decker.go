package decker

import (
	"image"
"github.com/corona10/goimagehash"
	"github.com/syndtr/goleveldb/leveldb"
)

// Decker is the main struct of the app
type Decker struct {
	Images    []image.Image
	DB        *leveldb.DB
	Threshold int
}

// Hash takes the perceptual hash of every image in the array and adds it it to the DB
func (d *Decker) Hash()  {
	for _, img := range d.Images {
		hash, err := goimagehash.PerceptualHash(img)
	}	
}
