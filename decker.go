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

type Graph struct {
	Threshold int
	Nodes     []Node
}

func NewGraph(threshold int) *Graph {
	return &Graph{
		Threshold: threshold,
		Nodes:     []Node{},
	}
}

func (g *Graph) Insert(img image.Image, hash *goimagehash.ImageHash, p string) (int, error) {
	node := Node{
		Image: img,
		Hash:  hash,
		Path:  p,
	}

	for i := range g.Nodes {
		distance, err := g.Nodes[i].Hash.Distance(hash)
		if err != nil {
			return -1, err
		}

		// Assume that there is only one possible candidate
		if distance <= g.Threshold {
			// If the new image is better
			if node.isBetterThan(&g.Nodes[i]) {
				// Steal the children
				node.Children = g.Nodes[i].Children
				// Insert current parent into children
				node.Children = append(node.Children, g.Nodes[i])
				// Become the parent
				g.Nodes[i] = node
				return i, nil
			}

			// Otherwise just append current one
			g.Nodes[i].Children = append(g.Nodes[i].Children, node)
			return i, nil
		}
	}

	// First time we've seen this image, make a new root
	g.Nodes = append(g.Nodes, node)
	return len(g.Nodes) - 1, nil
}

func (first *Node) isBetterThan(second *Node) bool {
	firstBounds := first.Image.Bounds()
	secondBounds := second.Image.Bounds()
	firstRes := firstBounds.Dx() * firstBounds.Dy()
	secondRes := secondBounds.Dx() * secondBounds.Dy()
	return firstRes > secondRes
}
