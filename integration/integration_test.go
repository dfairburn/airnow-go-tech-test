package integration

import (
	"github.com/dfairburn/airnow-go-tech-test/internal"
	"github.com/stretchr/testify/assert"
	"testing"
)

const testUrl = "https://crawler-test.com"

func Test_Walk_success(t *testing.T) {
	tree := internal.NewTree(testUrl)
	tree.Root.Walk(tree.Root, 0, 2)
	got := tree.Root.String("", 0)
	want := `https://crawler-test.com
  https://crawler-test.com//urls/double_slash/disallowed_start
  https://subdomain.crawler-test.com
`

	assert.Equal(t, got, want)
}

func TestWalk_earlyReturn(t *testing.T) {
	tree := internal.NewTree(testUrl)
	tree.Root.Walk(tree.Root, 1, 0)
	got := tree.Root.String("", 0)
	want := `https://crawler-test.com
`

	assert.Equal(t, got, want)
}
