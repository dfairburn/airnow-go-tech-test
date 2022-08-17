package internal

import (
	"fmt"
	"log"
	"net/url"
)

type Tree struct {
	Root *Node
}

type Node struct {
	url      string
	children []*Node
}

func NewTree(root string) *Tree {
	return &Tree{
		newNode(root),
	}
}

func newNode(url string) *Node {
	return &Node{
		url:      url,
		children: []*Node{},
	}
}

func (n *Node) Walk(root *Node, currLevel int, nestLevel int) {
	doc, err := get(n.url)
	if err != nil {
		root.remove(n.url)
		log.Println(err)
	}

	if currLevel >= nestLevel {
		return
	}

	children := getLinks(doc)
	for _, child := range children {
		n.insert(root, child)
	}

	for _, child := range n.children {
		child.Walk(root, currLevel+1, nestLevel)
	}
}

func (n *Node) String(s string, indentLevel int) string {
	indentation := ""
	for i := 0; i < indentLevel; i++ {
		indentation += "  "
	}

	s += fmt.Sprintf("%s%s\n", indentation, n.url)

	for _, child := range n.children {
		s = child.String(s, indentLevel+1)
	}

	return s
}

func (n *Node) remove(target string) {
	for i, child := range n.children {
		if child.url == target {
			newChildren := append(n.children[:i], n.children[i+1:]...)
			n.children = newChildren
			return
		}

		child.remove(target)
	}
}

func (n *Node) insert(root *Node, child *Node) {
	if !root.uniq(child.url) {
		return
	}

	n.children = append(n.children, child)
}

func (n *Node) uniq(raw string) bool {
	current := n.url

	if current == raw {
		return false
	}

	r, err := url.Parse(raw)
	if err != nil {
		return false
	}

	c, err := url.Parse(current)
	if err != nil {
		return false
	}

	if c.Host == r.Host && c.Path == r.Path {
		return false
	}

	if c.Path+"/" == r.Path || r.Path+"/" == c.Path {
		return false
	}

	for _, child := range n.children {
		return child.uniq(raw)
	}

	return true
}

func (n *Node) format(level int) string {
	var indentation string
	for i := 0; i < level; i++ {
		indentation = indentation + "  "
	}
	return fmt.Sprintf("%s%s\n", indentation, n.url)
}
