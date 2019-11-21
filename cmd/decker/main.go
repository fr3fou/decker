package main

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/fr3fou/decker"
)

func main() {
	if len(os.Args) < 1 {
		panic("please provide the directory")
	}

	dir := os.Args[1]

	d := decker.Decker{
		Threshold: 5,
		Mutex:     &sync.Mutex{},
	}

	err := filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		file, err := os.Open(p)
		if err != nil {
			return err
		}
		defer file.Close()

		ext := path.Ext(p)

		switch ext {
		case ".jpg", ".jpeg", ".png":
			img, fom, err := image.Decode(file)

			if err != nil {
				return err
			}

			log.Printf("%s encoded with format %s", path.Base(p), fom)

			go d.Hash(decker.Image{
				Image: img,
				Path:  p,
			})
		default:
			log.Printf("%s is an unsupported format %s", path.Base(p), ext)
			return nil
		}

		return nil
	})

	if err != nil {
		log.Println(err)
	}

	out, err := d.Check()

	if err != nil {
		panic(err)
	}

	dupeCount := 0
	for _, dupes := range out {
		for i := 1; i < len(dupes); i++ {
			log.Printf("%s is a duplicate of %s", path.Base(dupes[i].Path), path.Base(dupes[0].Path))
			dupeCount++
		}
	}

	log.Printf("%d dupes found", dupeCount)

}
