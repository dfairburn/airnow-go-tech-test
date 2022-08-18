package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/dfairburn/airnow-go-tech-test/internal"
	"time"
)

var target string
var timeout int
var nestLevel int

const (
	nestValue    = 0
	nestUsage    = "specified nesting level for traversal"
	targetValue  = ""
	targetUsage  = "the url to crawl"
	timeoutValue = 500
	timeoutUsage = "overall completion timeout in milliseconds"
)

func init() {
	flag.StringVar(&target, "target", targetValue, targetUsage)
	flag.StringVar(&target, "t", targetUsage, targetUsage+" (shorthand)")
	flag.IntVar(&nestLevel, "nest", nestValue, nestUsage)
	flag.IntVar(&nestLevel, "n", nestValue, nestUsage+" (shorthand)")
	flag.IntVar(&timeout, "timeout", timeoutValue, timeoutUsage)
	flag.IntVar(&timeout, "ti", timeoutValue, timeoutUsage+" (shorthand)")
}

func main() {
	flag.Parse()

	tree := internal.NewTree(target)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(timeout))
	defer cancel()

	tree.Root.Walk(ctx, tree.Root, 0, nestLevel)
	fmt.Println(tree.Root.String("", 0))
}
