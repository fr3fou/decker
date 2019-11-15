package decker

import (
	"encoding/binary"
	"fmt"
	"image"
	"log"

	"github.com/corona10/goimagehash"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
)

// Decker is the main struct of the app
type Decker struct {
	Images    map[string]image.Image
	DB        *leveldb.DB
	Threshold int
}

// Hash takes the perception hash of every image in the array and adds it it to the DB
func (d *Decker) Hash() {
	if d.DB == nil {
		panic("decker: leveldb.DB instance not provided")
	}

	for path, img := range d.Images {
		hash, err := goimagehash.PerceptionHash(img)
		if err != nil {
			log.Println(errors.Wrap(err, fmt.Sprintf("image %s couldn't be hashed", img)))
		}

		// Make the hash into a byte array
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, hash.GetHash())

		// Put into the DB
		d.DB.Put([]byte(path), b, nil)
	}
}

// Check checks all the images in the DB and returns the path for all duplicates
func (d *Decker) Check() ([]string, error) {
	if d.DB == nil {
		panic("decker: leveldb.DB instance not provided")
	}

	m := map[string][]byte{}
	iter := d.DB.NewIterator(nil, nil)

	for iter.Next() {
		// Get the key and val
		key := iter.Key()
		value := iter.Value()

		// Set the path and hash
		m[string(key)] = value
	}

	iter.Release()
	err := iter.Error()
	if err != nil {
		return []string{}, err
	}
}
