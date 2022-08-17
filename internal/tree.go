package internal

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
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
	doc, err := get(root, n.url)
	if err != nil {
		log.Println(err)
	}

	if currLevel == nestLevel {
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

func get(root *Node, target string) (*goquery.Document, error) {
	res, err := http.Get(target)
	if err != nil {
		root.remove(target)
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		root.remove(target)
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	return doc, err
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

func getLinks(doc *goquery.Document) []*Node {
	var children []*Node
	if doc == nil {
		return nil
	}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		rawUrl, ok := s.Attr("href")
		if ok {
			url, err := validUrl(rawUrl)
			if err != nil {
				return
			}

			n := newNode(url)
			children = append(children, n)
		}
	})

	return children
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

	for _, child := range n.children {
		return child.uniq(raw)
	}

	return true
}

func validUrl(raw string) (string, error) {
	var u *url.URL

	u, err := url.Parse(raw)
	if err == nil && u.Scheme != "" && u.Host != "" {
		return u.String(), nil
	}

	return "", fmt.Errorf("cannot parse url")
}

func (n *Node) format(level int) string {
	var indentation string
	for i := 0; i <= level; i++ {
		indentation = indentation + "  "
	}
	return fmt.Sprintf("%s%s", indentation, n.url)
}
