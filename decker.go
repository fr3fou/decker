package decker

import (
	"fmt"
	"image"
	"log"

	"github.com/corona10/goimagehash"
	"github.com/pkg/errors"
)

// Decker is the main struct of the app
type Decker struct {
	// TOOD: replace with channel?
	// somehow make concurrent
	// Output is a map of Hash -> []decker.Image
	Output map[uint64][]Image
	// Input is a map of Path -> image.Image
	Input map[string]image.Image
	// hashes is map that is used for
	// internal use cases (it holds the hash and ONE image)
	hashes map[uint64]Image
	// Threshold is the minimum hamming distance
	// for 2 images to be considered "different"
	Threshold int
}

// Hash takes the perception hash of every image in the array and makes a database
func (d *Decker) Hash() {
	for path, img := range d.Input {
		hash, err := goimagehash.PerceptionHash(img)

		if err != nil {
			log.Println(
				errors.Wrap(err,
					fmt.Sprintf("image %s couldn't be hashed", img),
				),
			)
		}
	}
}

// Check checks all the images in the DB and returns the path for all duplicates
func (d *Decker) Check() ([]Tree, error) {
	// // map of path -> hash
	// m := map[string][]byte{}
	// imgs := []Tree{}

	// iter.Release()
	// err := iter.Error()
	// if err != nil {
	// 	return []Tree{}, err
	// }

	// // Compare every hash with every hash *other* than it
	// for k1, v1 := range m {
	// 	for k2, v2 := range m {
	// 		// if we have the exact same path, carry on
	// 		if k2 == k1 {
	// 			continue
	// 		}

	// 		// convert the []byte to uint64
	// 		i1 := binary.LittleEndian.Uint64(v1)
	// 		i2 := binary.LittleEndian.Uint64(v2)

	// 		// get back the hashes
	// 		h1 := goimagehash.NewImageHash(i1, goimagehash.PHash)
	// 		h2 := goimagehash.NewImageHash(i2, goimagehash.PHash)

	// 		// calculate the hamming distance
	// 		distance, err := h1.Distance(h2)
	// 		if err != nil {
	// 			log.Println(
	// 				errors.Wrap(err,
	// 					fmt.Sprintf("decker: couldn't get the distance between %s and %s", k1, k2),
	// 				),
	// 			)
	// 		}

	// 		if distance <= d.Threshold {
	// 			// TOOD: find tree with coressponding hash
	// 			// update if it's a higher quality
	// 			// instead of appending, just mutate
	// 			tree := Tree{}
	// 			imgs = append(imgs, tree)
	// 		}
	// 	}
	// }

	return []Tree{}, nil
}
