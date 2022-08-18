// Package internal:
// Stores the definition of the Tree and node data structures and their methods
package internal

import (
	"context"
	"fmt"
	"net/url"
)

// Tree is the public data structure that is exposed from this package, as all operations are to be done on the root
// of the tree. The use of this data structure means we can easily separate out tiers/levels of links.
type Tree struct {
	Root *node
}

type node struct {
	url      string
	children []*node
}

// NewTree returns a new Tree struct initialized with a given root url to crawl
func NewTree(root string) *Tree {
	return &Tree{
		newNode(root),
	}
}

func newNode(url string) *node {
	return &node{
		url:      url,
		children: []*node{},
	}
}

// Walk queries the root url and recursively queries the children for their links until the depth is met.
// A call to get() is made first, as if this call fails, the node is removed from the tree so that it's cleared of
// broken urls as it continues to walk. It passes down a copy of the root node so that deeper levels in the tree
// can refer to the root not to remove children from and check uniqueness of the urls within the tree. It also doesn't
// allow insertion of the same link multiple times in the tree, although this doesn't currently handle the difference
// between relative links and absolute links.
func (n *node) Walk(ctx context.Context, root *node, currDepth int, depth int) {
	doc, err := get(ctx, n.url)
	if err != nil {
		root.remove(n.url)
		return
	}

	if currDepth >= depth {
		return
	}

	children := getLinks(root.url, doc)
	var curr string
	for _, child := range children {
		if curr == child.url {
			continue
		}
		n.insert(root, child)
		curr = child.url
	}

	for _, child := range n.children {
		child.Walk(ctx, root, currDepth+1, depth)
	}
}

// String returns the string representation of the tree. This method is a wrapper around the recursive indent method
// just to obfuscate having to call String() with initial string and indent values
func (n *node) String() string {
	return n.indent("", 0)
}

func (n *node) indent(s string, indentLevel int) string {
	indentation := ""
	for i := 0; i < indentLevel; i++ {
		indentation += "  "
	}

	s += fmt.Sprintf("%s%s\n", indentation, n.url)

	for _, child := range n.children {
		s = child.indent(s, indentLevel+1)
	}

	return s
}

// remove walks the tree from a given node to remove the first target node found that is matched on the node.url
func (n *node) remove(target string) {
	if len(n.children) == 0 {
		return
	}

	for i, child := range n.children {
		if child.url == target {
			newChildren := append(n.children[:i], n.children[i+1:]...)
			n.children = newChildren
		}

		child.remove(target)
	}
}

// insert checks if the child node to be inserted already appears in the tree. This means we're doing a O(n) traversal
// of the tree each time we insert, but I felt this was acceptable as it's an easier implementation to save time
func (n *node) insert(root *node, child *node) {
	if !root.uniq(child.url) {
		return
	}

	n.children = append(n.children, child)
}

// uniq finds if a given url is unique within a tree. This method decides that relative paths compared with absolute paths
// to the same endpoint are non-unique. It also decides that the difference between https and http is non-unique too.
// examples given a base url of http://foo-bar.com:
// http://foo-bar.com & https://foo-bar.com/ == false
// /baz & https://foo-bar.com/baz == false
func (n *node) uniq(raw string) bool {
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

// format is a helper method to indent the string format of nodes to the given level
func (n *node) format(level int) string {
	var indentation string
	for i := 0; i < level; i++ {
		indentation = indentation + "  "
	}
	return fmt.Sprintf("%s%s\n", indentation, n.url)
}
