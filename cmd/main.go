package main

import (
	"flag"
	"fmt"
	"github.com/dfairburn/airnow-go-tech-test/internal"
)

var target string
var nestLevel int

const (
	nestValue   = 0
	nestUsage   = "specified nesting level for traversal"
	targetValue = ""
	targetUsage = "the url to crawl"
)

func init() {
	flag.StringVar(&target, "target", targetValue, targetUsage)
	flag.StringVar(&target, "t", targetUsage, targetUsage+" (shorthand)")
	flag.IntVar(&nestLevel, "nest", nestValue, nestUsage)
	flag.IntVar(&nestLevel, "n", nestValue, nestUsage+" (shorthand)")
}

func main() {
	flag.Parse()

	tree := internal.NewTree(target)
	tree.Root.Walk(tree.Root, 0, nestLevel)
	fmt.Println(tree.Root.String("", 0))
}
