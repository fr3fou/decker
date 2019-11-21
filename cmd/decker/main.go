package main

import (
	"fmt"
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
	fmt.Println(runtime.NumCPU())
	ch := make(chan decker.Image, runtime.NumCPU())
	dir := os.Args[1]
	var wg sync.WaitGroup

	go func() {
		for k := range ch {
			imgs = append(imgs, k)
		}
	}()

	err := filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		wg.Add(1)
		defer wg.Done()

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

			go func() {
				ch <- decker.Hash(img, p)
			}()
		default:
			log.Printf("%s is an unsupported format %s", path.Base(p), ext)
			return nil
		}

		return nil
	})

	wg.Wait()
	close(ch)

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
