package decker

import (
	"fmt"
	"image"
	"log"

	"github.com/corona10/goimagehash"
	"github.com/pkg/errors"
)

type Output map[uint64][]Image

// Decker is the main struct of the app
type Decker struct {
	// TOOD: replace with channel?
	// somehow make concurrent
	// Output is a map of a randomly generated ID
	// ID -> []decker.Image
	// Input is a map of Path -> image.Image
	Input map[string]image.Image
	// hashes is map that is used for
	// internal use cases (it holds the hash and ONE image)
	// Hash -> Image
	hashes map[uint64]Image
	// Threshold is the minimum hamming distance
	// for 2 images to be considered "different"
	Threshold int
}

// Hash takes the perception hash of every image in the array and adds them to the hashes map
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

		key := hash.GetHash()

		// Add the hash
		d.hashes[key] = Image{
			img,
			path,
			hash,
			false,
		}
	}
}

// Check checks all the images in the DB and returns the path for all duplicates
func (d *Decker) Check() (Output, error) {
	output := Output{}

	for hash1, img1 := range d.hashes {
		for hash2, img2 := range d.hashes {
			// Ignore if we have the exact same image
			if hash1 == hash2 && img1.Path == img2.Path {
				continue
			}

			distance, err := h1.Distance(h2)
			if err != nil {
				log.Println(
					errors.Wrap(err,
						fmt.Sprintf("decker: couldn't get the distance between %s and %s", k1, k2),
					),
				)
			}

			// 	if distance <= d.Threshold {
			// 		// TOOD: find tree with coressponding hash
			// 		// update if it's a higher quality
			// 		// instead of appending, just mutate
			// 		tree := Tree{}
			// 		imgs = append(imgs, tree)
			// 	}
		}
	}

	return output
}
