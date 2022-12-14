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
var depth int

const (
	depthValue   = 0
	depthUsage   = "specified nesting level for traversal"
	targetValue  = ""
	targetUsage  = "https://crawler-test.com"
	timeoutValue = 5000
	timeoutUsage = "overall completion timeout in milliseconds"
)

func init() {
	flag.StringVar(&target, "target", targetValue, targetUsage)
	flag.StringVar(&target, "t", targetUsage, targetUsage+" (shorthand)")
	flag.IntVar(&depth, "depth", depthValue, depthUsage)
	flag.IntVar(&depth, "d", depthValue, depthUsage+" (shorthand)")
	flag.IntVar(&timeout, "timeout", timeoutValue, timeoutUsage)
	flag.IntVar(&timeout, "ti", timeoutValue, timeoutUsage+" (shorthand)")
}

func main() {
	flag.Parse()

	tree := internal.NewTree(target)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(timeout))
	defer cancel()

	fmt.Printf("crawling links from %s...\n", target)

	tree.Root.Walk(ctx, tree.Root, 0, depth)

	fmt.Println(tree.Root.String())
}
