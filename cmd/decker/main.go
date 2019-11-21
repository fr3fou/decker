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

	imgs := []decker.Image{}

	dir := os.Args[1]
	m := &sync.Mutex{}

	var wg sync.WaitGroup

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
			wg.Add(1)
			img, fom, err := image.Decode(file)

			if err != nil {
				wg.Done()
				return err
			}

			log.Printf("%s encoded with format %s", path.Base(p), fom)

			go func() {
				i := decker.Hash(img, p)

				m.Lock()
				imgs = append(imgs, i)
				m.Unlock()

				wg.Done()
			}()
		default:
			log.Printf("%s is an unsupported format %s", path.Base(p), ext)
			return nil
		}

		return nil
	})

	wg.Wait()

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
