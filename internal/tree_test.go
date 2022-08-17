package internal

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const testUrl = "https://crawler-test.com/"
const invalidUrl = "invalid"

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

func Test_get(t *testing.T) {
	type args struct {
		root   *Node
		target string
	}
	tests := []struct {
		name string
		args args
		err  assert.ErrorAssertionFunc
		want assert.ValueAssertionFunc
	}{
		{
			name: "given a valid test url doesn't return an error",
			args: args{
				root: &Node{
					url:      testUrl,
					children: []*Node{},
				},
				target: testUrl,
			},
			err:  assert.NoError,
			want: assert.NotNil,
		},
		{
			name: "given an invalid test url return an error",
			args: args{
				root: &Node{
					url:      testUrl,
					children: []*Node{},
				},
				target: invalidUrl,
			},
			err:  assert.Error,
			want: assert.Nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := get(tt.args.root, tt.args.target)
			tt.want(t, got)
			tt.err(t, err)
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Node{
				url:      tt.fields.url,
				children: tt.fields.children,
			}
			got := n.uniq(tt.args.raw)

			assert.Equalf(t, tt.want, got, "uniq(%v)", tt.args.raw)
		})
	}
}

func Test_validUrl(t *testing.T) {
	type args struct {
		raw string
	}
	tests := []struct {
		name string
		args args
		want string
		err  assert.ErrorAssertionFunc
	}{
		{
			name: "returns url given a valid url",
			args: args{
				raw: testUrl,
			},
			want: testUrl,
			err:  assert.NoError,
		},
		{
			name: "returns error given an invalid url",
			args: args{
				raw: invalidUrl,
			},
			want: "",
			err:  assert.Error,
		},
		{
			name: "returns error given a path",
			args: args{
				raw: "/path",
			},
			want: "",
			err:  assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validUrl(tt.args.raw)
			if !tt.err(t, err, fmt.Sprintf("validUrl(%v)", tt.args.raw)) {
				return
			}
			assert.Equalf(t, tt.want, got, "validUrl(%v)", tt.args.raw)
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
