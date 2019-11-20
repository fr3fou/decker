package decker

import (
	"fmt"
	"image"
	"log"

	"github.com/corona10/goimagehash"
	"github.com/pkg/errors"
)

// Output is a map of a generated ID and decker.Image
// ID -> []decker.Image
type Output map[uint64][]Image

// Decker is the main struct of the app
type Decker struct {
	// TOOD: replace with channel?
	// somehow make concurrent
	// Input is a map of Path -> image.Image
	Input map[string]image.Image
	// hashes is an array for
	// internal use cases
	hashes []Image
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

		// Add the hash
		d.hashes = append(d.hashes, Image{
			Image:  img,
			Path:   path,
			ID:     0,
			Hash:   hash,
			IsBest: false,
		})
	}
}

// Check checks all the images in the DB and returns the output
func (d *Decker) Check() (Output, error) {
	output := Output{}

	// when making concurrent, use a mutex or a random UUID?
	var id uint64 = 0

	// Compare each image with eachother
	for _, img1 := range d.hashes {
		// if it's an image that we have seen before
		// (already exists in our map)
		// we should just carry on
		if _, ok := output[img1.ID]; ok {
			continue
		}

		// Update the ID if it's a new image
		id++
		img1.ID = id

		// Add the image first
		// Can this be optimized?
		// img1 might not have any duplicates so we are wasting space
		output[img1.ID] = []Image{}
		output[img1.ID] = append(output[img1.ID], img1)

		// Compare to the rest of the images
		// We have to use a C-Style for loop, because we are going to be mutating
		for i := 0; i < len(d.hashes); i++ {
			img2 := d.hashes[i]
			// Ignore if we have the exact same image
			if img1.Path == img2.Path {
				continue
			}

			// Get the actual hashes
			h1, h2 := img1.Hash, img2.Hash

			// Calculate the hamming distance
			distance, err := h1.Distance(h2)
			if err != nil {
				log.Println(
					errors.Wrap(err,
						fmt.Sprintf("decker: couldn't calculate the distance between %s and %s", img1.Path, img2.Path),
					),
				)
				continue
			}

			if distance <= d.Threshold {
				// If the images are duplicates
				img2.ID = img1.ID
				d.hashes[i].ID = img1.ID

				// Add the current image
				output[img1.ID] = append(output[img1.ID], img2)
			}
		}

		// If we have no duplicates
		if len(output[img1.ID]) == 1 {
			// Delete the entry
			delete(output, img1.ID)
		}
	}

	return output, nil
}
