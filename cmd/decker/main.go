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

	"github.com/fr3fou/decker"
)

func main() {
	if len(os.Args) < 1 {
		panic("please provide the directory")
	}

	dir := os.Args[1]

	var imgs []image.Image

	err := filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		file, err := os.Open(p)
		if err != nil {
			return err
		}

		ext := path.Ext(p)

		switch ext {
		case ".jpg", ".jpeg", ".png":
			img, fom, err := image.Decode(file)

			if err != nil {
				return err
			}

			log.Printf("%s encoded with format %s", path.Base(p), fom)

			imgs = append(imgs, img)
		default:
			log.Printf("%s is an unsupported format %s", path.Base(p), ext)
			return nil
		}

		return nil
	})

	if err != nil {
		log.Println(err)
	}

	_ = decker.Decker{}
}
