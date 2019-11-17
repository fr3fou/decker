package decker

import (
	"fmt"
	"image"
	"log"

	"github.com/corona10/goimagehash"
	"github.com/pkg/errors"
)

// Output is a map of a randomly generated ID and decker.Image
// ID -> []decker.Image
type Output map[uint64][]*Image

// Decker is the main struct of the app
type Decker struct {
	// TOOD: replace with channel?
	// somehow make concurrent
	// Input is a map of Path -> image.Image
	Input map[string]image.Image
	// hashes is an array for
	// internal use cases
	hashes []*Image
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
					fmt.Sprintf("decker: image %s couldn't be hashed", img),
				),
			)
		}

		key := hash.GetHash()

		// Add the hash
		d.hashes = append(d.hashes, &Image{
			img,
			path,
			-1,
			hash,
			false,
		})
	}
}

// Check checks all the images in the DB and returns the output
func (d *Decker) Check() (Output, error) {
	output := Output{}

	// when making concurrent, use a mutex or a random UUID?
	id := 0

	for _, img1 := range d.hashes {
		// set the ID if it's an image that we have not seen before
		if img1.ID == -1 {
			id++
			img1.ID = id
		}

		for _, img2 := range d.hashes {
			// Ignore if we have the exact same image
			if img1.Path == img2.Path {
				continue
			}

			h1, h2 := img1.Hash, img2.Hash

			distance, err := h1.Distance(h2)
			if err != nil {
				log.Println(
					errors.Wrap(err,
						fmt.Sprintf("decker: couldn't get the distance between %s and %s", img1.Path, img2.Path),
					),
				)
			}

			if distance <= d.Threshold {
				img2.ID = img1.ID

				if key, ok := output[id]; !ok {
					output[id] = []*Image{}
				}

				output[id] = append(output, img2)
			}
		}
	}

	return output, nil
}
