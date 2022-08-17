package internal

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	testUrl     = "https://crawler-test.com/"
	redirectUrl = "https://crawler-test.com/redirects/redirect_to_404"
	invalidUrl  = "invalid"
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
				Root: &Node{
					url:      testUrl,
					children: []*Node{},
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
		children []*Node
	}
	type args struct {
		target string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Node
	}{
		{
			name: "removes a given element from the tree",
			fields: fields{
				url: "a",
				children: []*Node{
					{
						url:      "b",
						children: []*Node{},
					},
					{
						url: "c",
						children: []*Node{
							{
								url:      "e",
								children: []*Node{},
							},
						},
					},
					{
						url:      "d",
						children: []*Node{},
					},
				},
			},
			args: args{
				target: "e",
			},
			want: &Node{
				url: "a",
				children: []*Node{
					{
						url:      "b",
						children: []*Node{},
					},
					{
						url: "c",
						children: []*Node{
							{
								url:      "e",
								children: []*Node{},
							},
						},
					},
					{
						url:      "d",
						children: []*Node{},
					},
				},
			},
		},
		{
			name: "doesn't remove element that doesn't exist",
			fields: fields{
				url: "a",
				children: []*Node{
					{
						url:      "b",
						children: []*Node{},
					},
				},
			},
			args: args{
				target: "e",
			},
			want: &Node{
				url: "a",
				children: []*Node{
					{
						url:      "b",
						children: []*Node{},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Node{
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
		children []*Node
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
				children: []*Node{
					{
						url:      "b",
						children: []*Node{},
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
				children: []*Node{
					{
						url:      "b",
						children: []*Node{},
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
				children: []*Node{
					{
						url:      "b",
						children: []*Node{},
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
				children: []*Node{},
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
				children: []*Node{},
			},
			args: args{
				raw: "https://a_url.com/",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Node{
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
		children []*Node
	}
	type args struct {
		child *Node
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Node
	}{
		{
			name: "inserts a child into the root node",
			fields: fields{
				url:      "a",
				children: []*Node{},
			},
			args: args{
				child: &Node{
					url:      "b",
					children: []*Node{},
				},
			},
			want: &Node{
				url: "a",
				children: []*Node{
					{
						url:      "b",
						children: []*Node{},
					},
				},
			},
		},
		{
			name: "doesn't insert if child already exists in the tree",
			fields: fields{
				url: "a",
				children: []*Node{
					{
						url:      "b",
						children: []*Node{},
					},
				},
			},
			args: args{
				child: &Node{
					url:      "b",
					children: []*Node{},
				},
			},
			want: &Node{
				url: "a",
				children: []*Node{
					{
						url:      "b",
						children: []*Node{},
					},
				},
			},
		},
		{
			name: "doesn't insert if host and path is the same",
			fields: fields{
				url: "http://a",
				children: []*Node{
					{
						url:      "http://b",
						children: []*Node{},
					},
				},
			},
			args: args{
				child: &Node{
					url:      "https://a",
					children: []*Node{},
				},
			},
			want: &Node{
				url: "http://a",
				children: []*Node{
					{
						url:      "http://b",
						children: []*Node{},
					},
				},
			},
		},
		{
			name: "inserts if host are same but paths are different",
			fields: fields{
				url: "http://a.com",
				children: []*Node{
					{
						url:      "http://b.com",
						children: []*Node{},
					},
				},
			},
			args: args{
				child: &Node{
					url:      "https://a.com/b/c",
					children: []*Node{},
				},
			},
			want: &Node{
				url: "http://a.com",
				children: []*Node{
					{
						url:      "http://b.com",
						children: []*Node{},
					},
					{
						url:      "https://a.com/b/c",
						children: []*Node{},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Node{
				url:      tt.fields.url,
				children: tt.fields.children,
			}
			n.insert(n, tt.args.child)

			if n.diff(tt.want) {
				t.Errorf("wrong object recieved\ngot=%s\nwant=%s", n.String("", 0), tt.want.String("", 0))
			}
		})
	}
}

func (n *Node) diff(m *Node) bool {
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

func TestNode_format(t *testing.T) {
	type fields struct {
		url      string
		children []*Node
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
				children: []*Node{},
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
				children: []*Node{},
			},
			args: args{
				level: 2,
			},
			want: fmt.Sprintf("    a\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Node{
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
		children []*Node
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
				children: []*Node{
					{
						url: "b",
						children: []*Node{
							{
								url:      "c",
								children: []*Node{},
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
			n := &Node{
				url:      tt.fields.url,
				children: tt.fields.children,
			}
			assert.Equalf(t, tt.want, n.String(tt.args.s, tt.args.indentLevel), "String(%v, %v)", tt.args.s, tt.args.indentLevel)
		})
	}
}
