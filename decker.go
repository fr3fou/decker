package decker

// Decker is the main struct of the app
type Decker struct {
	Images    []Tree
	Threshold int
}

// Hash takes the perception hash of every image in the array and makes a database
func (d *Decker) Hash() {
	// for path, img := range d.Images {
	// 	hash, err := goimagehash.PerceptionHash(img)
	// 	if err != nil {
	// 		log.Println(
	// 			errors.Wrap(err,
	// 				fmt.Sprintf("image %s couldn't be hashed", img),
	// 			),
	// 		)
	// 	}

	// 	// Make the hash into a byte array
	// 	b := make([]byte, 8)
	// 	binary.LittleEndian.PutUint64(b, hash.GetHash())

	// 	// Put into the DB
	// 	// d.DB.Put([]byte(path), b, nil)
	// }
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

	// return imgs, nil
}
