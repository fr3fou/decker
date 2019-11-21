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
	"runtime"
	"sync"

	"github.com/fr3fou/decker"
)

func main() {
	if len(os.Args) < 1 {
		panic("please provide the directory")
	}

	imgs := []decker.Image{}

	dir := os.Args[1]
	m := &sync.Mutex{}
	sem := make(chan int, runtime.NumCPU())

	err := filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		ext := path.Ext(p)

		switch ext {
		case ".jpg", ".jpeg", ".png":
			sem <- 1

			file, err := os.Open(p)
			if err != nil {
				return err
			}

			go func() {
				defer func() { <-sem }()
				defer file.Close()

				img, fom, err := image.Decode(file)

				if err != nil {
					log.Printf("couldn't decode %s", path.Base(p))
					return
				}

				log.Printf("%s decoded with format %s", path.Base(p), fom)

				i := decker.Hash(img, p)

				m.Lock()
				imgs = append(imgs, i)
				m.Unlock()
			}()
		default:
			log.Printf("%s is an unsupported format %s", path.Base(p), ext)
			return nil
		}

		return nil
	})

	for i := 0; i < runtime.NumCPU(); i++ {
		sem <- 1
	}

	if err != nil {
		log.Println(err)
	}

	out, err := decker.Check(imgs, 5)

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
