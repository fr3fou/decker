package decker

import (
	"fmt"
	"image"
	"log"
	"path"

	"github.com/corona10/goimagehash"
	"github.com/pkg/errors"
)

// Output is a map of a generated ID and decker.Image
// ID -> []decker.Image
type Output map[uint64][]Image

// Hash takes the perception hash of an image and returns it
func Hash(img image.Image, p string) (*Image, error) {
	hash, err := goimagehash.PerceptionHash(img)

	if err != nil {
		return nil, errors.Wrap(err,
			fmt.Sprintf("decker: image %s couldn't be hashed", p),
		)
	}

	log.Printf("%s hashed with hash %x", path.Base(p), hash.GetHash())

	return &Image{
		Image:  img,
		Hash:   hash,
		ID:     0,
		IsBest: false,
		Path:   p,
	}, nil
}

// Check checks all the images in the DB and returns the output
func Check(hashes []Image, threshold int) (Output, error) {
	output := Output{}

	// when making concurrent, use a mutex or a random UUID?
	var id uint64 = 0

	// Compare each image with eachother
	for _, img1 := range hashes {
		// if it's an image that we have seen before
		// (already exists in our map)
		// we should just carry on
		if _, ok := output[img1.ID]; ok {
			continue
		}

		// Update the ID if it's a new image
		id++
		img1.ID = id

		output[img1.ID] = []Image{}
		output[img1.ID] = append(output[img1.ID], img1)

		// Compare to the rest of the images
		// We have to use a C-Style for loop, because we are going to be mutating
		for i := 0; i < len(hashes); i++ {
			img2 := hashes[i]
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

			if distance <= threshold {
				// If the images are duplicates
				img2.ID = img1.ID
				hashes[i].ID = img1.ID

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

	checkForBest(output)

	return output, nil
}

func checkForBest(imgs Output) {
	for _, v := range imgs {
		bestIdx := -1
		bounds1 := v[0].Bounds()
		res1 := bounds1.Dx() * bounds1.Dy()
		for i := 1; i < len(v); i++ {
			bounds2 := v[i].Bounds()
			res2 := bounds2.Dx() * bounds2.Dy()
			if res2 > res1 {
				if bestIdx > 0 {
					// Update the previous one
					v[bestIdx].IsBest = false
				}

				bestIdx = i

				v[bestIdx].IsBest = true

				res1 = res2
				bounds1 = bounds2
			}
		}

		if bestIdx > 0 {
			// Swap
			v[0], v[bestIdx] = v[bestIdx], v[0]
		}
	}
}
