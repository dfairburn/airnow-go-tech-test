package integration

import (
	"context"
	"github.com/dfairburn/airnow-go-tech-test/internal"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const testUrl = "https://crawler-test.com/links/relative_link/a/b"
const timeout = 1000

func Test_Walk_success(t *testing.T) {
	tree := internal.NewTree(testUrl)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()

	tree.Root.Walk(ctx, tree.Root, 0, 2)
	got := tree.Root.String()
	want := `https://crawler-test.com/links/relative_link/a/b
  https://crawler-test.com/links/relative_link/a/by/z
`

	assert.Equal(t, got, want)
}

func TestWalk_earlyReturn(t *testing.T) {
	tree := internal.NewTree(testUrl)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()

	tree.Root.Walk(ctx, tree.Root, 1, 0)
	got := tree.Root.String()
	want := `https://crawler-test.com/links/relative_link/a/b
`

	assert.Equal(t, got, want)
}
