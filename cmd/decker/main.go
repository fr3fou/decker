package main

import (
	"flag"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/fr3fou/decker"
)

func main() {
	dir := ""
	flag.StringVar(&dir, "dir", "", "path to the directory which contains the images")
	flag.StringVar(&dir, "d", "", "path to the directory which contains the images")

	threshold := 5
	flag.IntVar(&threshold, "threshold", 5, "threshold amount")
	flag.IntVar(&threshold, "t", 5, "threshold amount")

	flag.Parse()

	if dir == "" {
		panic("dir flag is required")
	}

	imgs := []decker.Image{}

	// Mutex for writing to the imgs array
	m := &sync.Mutex{}

	// Semaphore due to `ulimit`
	sem := make(chan int, runtime.NumCPU())

	err := filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		ext := strings.ToLower(path.Ext(p))

		switch ext {
		case ".jpg", ".jpeg", ".png":
			// Block here, as there's a limited amount of files open at a given time
			// Check `ulimit -n`
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

	// Add the last jobs
	for i := 0; i < runtime.NumCPU(); i++ {
		sem <- 1
	}

	if err != nil {
		log.Println(err)
	}

	out, err := decker.Check(imgs, threshold)

	if err != nil {
		panic(err)
	}

	dupeCount := 0
	for _, dupes := range out {
		log.Printf("%s has a resolution of %dx%d and has %d dupes", path.Base(dupes[0].Path), dupes[0].Bounds().Dx(), dupes[0].Bounds().Dy(), len(dupes)-1)
		for i := 1; i < len(dupes); i++ {
			log.Printf("\t%s is a duplicate of %s", path.Base(dupes[i].Path), path.Base(dupes[0].Path))
			dupeCount++
		}
	}

	log.Printf("%d dupes found", dupeCount)
}
