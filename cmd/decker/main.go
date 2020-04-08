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

	"github.com/corona10/goimagehash"
	"github.com/fr3fou/decker"
)

type finishedEvent struct{}

func main() {
	dir := ""
	flag.StringVar(&dir, "dir", "", "path to the directory which contains the images")
	flag.StringVar(&dir, "d", "", "path to the directory which contains the images")

	threshold := 5
	flag.IntVar(&threshold, "threshold", 5, "threshold amount")
	flag.IntVar(&threshold, "t", 5, "threshold amount")

	flag.Parse()

	if dir == "" {
		var err error
		dir, err = os.Getwd()

		if err != nil {
			panic(err)
		}
	}

	tree := decker.NewTree(threshold)

	// Mutex for writing to the imgs array
	m := &sync.Mutex{}

	// Semaphore due to `ulimit`
	sem := make(chan finishedEvent, runtime.NumCPU())

	err := filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		ext := strings.ToLower(path.Ext(p))

		switch ext {
		case ".jpg", ".jpeg", ".png":
			// Block here, as there's a limited amount of files open at a given time
			// Check `ulimit -n`
			sem <- finishedEvent{}

			file, err := os.Open(p)
			if err != nil {
				return err
			}

			go func() {
				defer func() { <-sem }()

				img, fom, err := image.Decode(file)
				file.Close()

				if err != nil {
					log.Printf("couldn't decode %s", path.Base(p))
					return
				}

				log.Printf("%s decoded with format %s", path.Base(p), fom)

				hash, err := goimagehash.PerceptionHash(img)
				if err != nil {
					log.Printf("couldn't hash image %s", p)
					return
				}
				m.Lock()
				_, err = tree.Insert(img, hash, p)
				m.Unlock()
				if err != nil {
					log.Printf("couldn't insert image %s into tree", p)
					return
				}
			}()
		default:
			log.Printf("%s is an unsupported format %s", path.Base(p), ext)
			return nil
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	// Add the last jobs
	for i := 0; i < runtime.NumCPU(); i++ {
		sem <- finishedEvent{}
	}

	dupeCount := 0
	for _, node := range tree.Nodes {
		dupes := node.Children
		if len(dupes) == 0 {
			continue
		}
		log.Printf("%s has a resolution of %dx%d and has %d dupes", path.Base(node.Path), node.Image.Bounds().Dx(), node.Image.Bounds().Dy(), len(dupes))
		for i := 0; i < len(dupes); i++ {
			log.Printf("\t%s is a duplicate of %s", path.Base(dupes[i].Path), path.Base(node.Path))
			dupeCount++
		}
	}

	log.Printf("%d dupes found", dupeCount)
}
