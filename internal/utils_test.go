package internal

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_get(t *testing.T) {
	type args struct {
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
				target: testUrl,
			},
			err:  assert.NoError,
			want: assert.NotNil,
		},
		{
			name: "given an invalid test url return an error",
			args: args{
				target: invalidUrl,
			},
			err:  assert.Error,
			want: assert.Nil,
		},
		{
			name: "given an invalid test url return an error",
			args: args{
				target: brokenUrl,
			},
			err:  assert.Error,
			want: assert.Nil,
		},
		{
			name: "given a url that redirects returns an error",
			args: args{
				target: redirectUrl,
			},
			err:  assert.Error,
			want: assert.Nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
			defer cancel()

			got, err := get(ctx, tt.args.target)
			tt.want(t, got)
			tt.err(t, err)
		})
	}
}

func Test_validUrl(t *testing.T) {
	type args struct {
		raw  string
		base string
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
				raw:  testUrl,
				base: testUrl,
			},
			want: testUrl,
			err:  assert.NoError,
		},
		{
			name: "returns error given an invalid url",
			args: args{
				raw:  invalidUrl,
				base: testUrl,
			},
			want: "",
			err:  assert.Error,
		},
		{
			name: "returns full url given a relative path",
			args: args{
				raw:  "/path",
				base: baseUrl,
			},
			want: "https://crawler-test.com/path",
			err:  assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validUrl(tt.args.base, tt.args.raw)
			if !tt.err(t, err, fmt.Sprintf("validUrl(%v)", tt.args.raw)) {
				return
			}
			assert.Equalf(t, tt.want, got, "validUrl(%v)", tt.args.raw)
		})
	}
}

func Test_getLinks(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		want *node
		args args
	}{
		{
			name: "get links retrieves all hrefs on a page",
			args: args{
				url: testUrl,
			},
			want: &node{
				url: "https://crawler-test.com/links/relative_link/a/b",
				children: []*node{
					{
						url:      "https://crawler-test.com/links/relative_link/a/b/",
						children: []*node{},
					},
					{
						url:      "https://crawler-test.com/links/relative_link/a/by/z",
						children: []*node{},
					},
					{
						url:      "https://crawler-test.com/links/relative_link/a/b/content/custom_text/relative_link_with_a_slash_at_the_beginning_target",
						children: []*node{},
					},
					{
						url:      "https://crawler-test.com/links/relative_link/a/b?parameter_only_link=1",
						children: []*node{},
					},
				},
			},
		},
		{
			name: "returns early with an invalid url",
			args: args{
				url: notFoundUrl,
			},
			want: &node{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := NewTree(tt.args.url)
			doc, err := get(context.Background(), root.Root.url)
			if err != nil {
				return
			}

			root.Root.children = getLinks(tt.args.url, doc)
			if root.Root.diff(tt.want) {
				t.Errorf("wrong object recieved, got=%v\nwant=%v", root.Root.String(), tt.want.String())
			}
		})
	}
}
