package internal

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	baseUrl     = "https://crawler-test.com"
	testUrl     = "https://crawler-test.com/links/relative_link/a/b"
	redirectUrl = "https://crawler-test.com/redirects/redirect_to_404"
	notFoundUrl = "http://crawler-test.com/status_codes/status_400"
	brokenUrl   = "https://invalid.crawler-test.com/"
	invalidUrl  = "::invalid::"
	timeout     = 1000
)

func TestNewTree(t *testing.T) {
	type args struct {
		root string
	}
	tests := []struct {
		name string
		args args
		want *Tree
	}{
		{
			name: "given a url string input",
			args: args{
				root: testUrl,
			},
			want: &Tree{
				Root: &node{
					url:      testUrl,
					children: []*node{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTree(tt.args.root)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestNode_remove(t *testing.T) {
	type fields struct {
		url      string
		children []*node
	}
	type args struct {
		target string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *node
	}{
		{
			name: "removes a given element from the tree",
			fields: fields{
				url: "a",
				children: []*node{
					{
						url:      "b",
						children: []*node{},
					},
					{
						url: "c",
						children: []*node{
							{
								url:      "e",
								children: []*node{},
							},
						},
					},
					{
						url:      "d",
						children: []*node{},
					},
				},
			},
			args: args{
				target: "e",
			},
			want: &node{
				url: "a",
				children: []*node{
					{
						url:      "b",
						children: []*node{},
					},
					{
						url: "c",
						children: []*node{
							{
								url:      "e",
								children: []*node{},
							},
						},
					},
					{
						url:      "d",
						children: []*node{},
					},
				},
			},
		},
		{
			name: "doesn't remove element that doesn't exist",
			fields: fields{
				url: "a",
				children: []*node{
					{
						url:      "b",
						children: []*node{},
					},
				},
			},
			args: args{
				target: "e",
			},
			want: &node{
				url: "a",
				children: []*node{
					{
						url:      "b",
						children: []*node{},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &node{
				url:      tt.fields.url,
				children: tt.fields.children,
			}
			n.remove(tt.args.target)

			if n.diff(tt.want) {
				t.Errorf("wrong object recieved, got=%v\nwant=%v", n, tt.want)
			}

		})
	}
}

func TestNode_uniq(t *testing.T) {
	type fields struct {
		url      string
		children []*node
	}
	type args struct {
		raw string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "returns true that a given element is unique in the tree",
			fields: fields{
				url: "a",
				children: []*node{
					{
						url:      "b",
						children: []*node{},
					},
				},
			},
			args: args{
				raw: "c",
			},
			want: true,
		},
		{
			name: "returns false that a given element already exists",
			fields: fields{
				url: "a",
				children: []*node{
					{
						url:      "b",
						children: []*node{},
					},
				},
			},
			args: args{
				raw: "b",
			},
			want: false,
		},
		{
			name: "returns false if given url is invalid",
			fields: fields{
				url: "a",
				children: []*node{
					{
						url:      "b",
						children: []*node{},
					},
				},
			},
			args: args{
				raw: "::invalid::",
			},
			want: false,
		},
		{
			name: "returns false if node url is invalid",
			fields: fields{
				url:      "::invalid::",
				children: []*node{},
			},
			args: args{
				raw: "a",
			},
			want: false,
		},
		{
			name: "returns false if urls only differ by trailing '/'",
			fields: fields{
				url:      "https://a_url.com",
				children: []*node{},
			},
			args: args{
				raw: "https://a_url.com/",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &node{
				url:      tt.fields.url,
				children: tt.fields.children,
			}
			got := n.uniq(tt.args.raw)

			assert.Equalf(t, got, tt.want, "uniq(%v)", tt.args.raw)
		})
	}
}

func TestNode_insert(t *testing.T) {
	type fields struct {
		url      string
		children []*node
	}
	type args struct {
		child *node
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *node
	}{
		{
			name: "inserts a child into the root node",
			fields: fields{
				url:      "a",
				children: []*node{},
			},
			args: args{
				child: &node{
					url:      "b",
					children: []*node{},
				},
			},
			want: &node{
				url: "a",
				children: []*node{
					{
						url:      "b",
						children: []*node{},
					},
				},
			},
		},
		{
			name: "doesn't insert if child already exists in the tree",
			fields: fields{
				url: "a",
				children: []*node{
					{
						url:      "b",
						children: []*node{},
					},
				},
			},
			args: args{
				child: &node{
					url:      "b",
					children: []*node{},
				},
			},
			want: &node{
				url: "a",
				children: []*node{
					{
						url:      "b",
						children: []*node{},
					},
				},
			},
		},
		{
			name: "doesn't insert if host and path is the same",
			fields: fields{
				url: "http://a",
				children: []*node{
					{
						url:      "http://b",
						children: []*node{},
					},
				},
			},
			args: args{
				child: &node{
					url:      "https://a",
					children: []*node{},
				},
			},
			want: &node{
				url: "http://a",
				children: []*node{
					{
						url:      "http://b",
						children: []*node{},
					},
				},
			},
		},
		{
			name: "inserts if host are same but paths are different",
			fields: fields{
				url: "http://a.com",
				children: []*node{
					{
						url:      "http://b.com",
						children: []*node{},
					},
				},
			},
			args: args{
				child: &node{
					url:      "https://a.com/b/c",
					children: []*node{},
				},
			},
			want: &node{
				url: "http://a.com",
				children: []*node{
					{
						url:      "http://b.com",
						children: []*node{},
					},
					{
						url:      "https://a.com/b/c",
						children: []*node{},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &node{
				url:      tt.fields.url,
				children: tt.fields.children,
			}
			n.insert(n, tt.args.child)

			if n.diff(tt.want) {
				t.Errorf("wrong object recieved\ngot=%s\nwant=%s", n.String(), tt.want.String())
			}
		})
	}
}

func TestNode_format(t *testing.T) {
	type fields struct {
		url      string
		children []*node
	}
	type args struct {
		level int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "formats the node with no indentation",
			fields: fields{
				url:      "a",
				children: []*node{},
			},
			args: args{
				level: 0,
			},
			want: fmt.Sprintf("a\n"),
		},
		{
			name: "formats the node with indentation",
			fields: fields{
				url:      "a",
				children: []*node{},
			},
			args: args{
				level: 2,
			},
			want: fmt.Sprintf("    a\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &node{
				url:      tt.fields.url,
				children: tt.fields.children,
			}
			assert.Equalf(t, tt.want, n.format(tt.args.level), "format(%v)", tt.args.level)
		})
	}
}

func TestNode_String(t *testing.T) {
	type fields struct {
		url      string
		children []*node
	}
	type args struct {
		s           string
		indentLevel int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "gives the string representation of the tree",
			fields: fields{
				url: "a",
				children: []*node{
					{
						url: "b",
						children: []*node{
							{
								url:      "c",
								children: []*node{},
							},
						},
					},
				},
			},
			args: args{
				s:           "",
				indentLevel: 0,
			},
			want: "a\n  b\n    c\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &node{
				url:      tt.fields.url,
				children: tt.fields.children,
			}
			assert.Equalf(t, tt.want, n.String(), "String(%v, %v)", tt.args.s, tt.args.indentLevel)
		})
	}
}

func (n *node) diff(m *node) bool {
	if n.url != m.url {
		fmt.Println(n.url, m.url)
		return true
	}

	if len(n.children) != len(m.children) {
		return true
	}

	for i, child := range n.children {
		return child.diff(m.children[i])
	}

	return false
}
