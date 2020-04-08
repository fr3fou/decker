package decker

import (
	"image"

	"github.com/corona10/goimagehash"
)

type Node struct {
	Image    image.Image
	Path     string
	Hash     *goimagehash.ImageHash
	Children []Node
}

type Tree struct {
	Threshold int
	Nodes     []Node
}

func NewTree(threshold int) *Tree {
	return &Tree{
		Threshold: threshold,
		Nodes:     []Node{},
	}
}

func (t *Tree) Insert(img image.Image, hash *goimagehash.ImageHash, p string) (int, error) {
	node := Node{
		Image: img,
		Hash:  hash,
		Path:  p,
	}

	for i := range t.Nodes {
		distance, err := t.Nodes[i].Hash.Distance(hash)
		if err != nil {
			return -1, err
		}

		// Assume that there is only one possible candidate
		if distance <= t.Threshold {
			// If the new image is better
			if node.isBetterThan(&t.Nodes[i]) {
				// Steal the children
				node.Children = t.Nodes[i].Children
				// Insert current parent into children
				node.Children = append(node.Children, t.Nodes[i])
				// Become the parent
				t.Nodes[i] = node
				return i, nil
			}

			// Otherwise just append current one
			t.Nodes[i].Children = append(t.Nodes[i].Children, node)
			return i, nil
		}
	}

	// First time we've seen this image, make a new root
	t.Nodes = append(t.Nodes, node)
	return len(t.Nodes) - 1, nil
}

func (first *Node) isBetterThan(second *Node) bool {
	firstBounds := first.Image.Bounds()
	secondBounds := second.Image.Bounds()
	firstRes := firstBounds.Dx() * firstBounds.Dy()
	secondRes := secondBounds.Dx() * secondBounds.Dy()
	return firstRes > secondRes
}
