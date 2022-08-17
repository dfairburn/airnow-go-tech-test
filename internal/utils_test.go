package internal

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
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
			got, err := get(tt.args.target)
			tt.want(t, got)
			tt.err(t, err)
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

func Test_getLinks(t *testing.T) {
	tests := []struct {
		name string
		want []*Node
		err  assert.ErrorAssertionFunc
	}{
		{
			name: "get links retrieves all hrefs on a page",
			want: []*Node{
				{
					url:      "http://www.s%C3%B8kbar.no",
					children: []*Node{},
				},
				{
					url:      "https://crawler-test.com//urls/double_slash/disallowed_start",
					children: []*Node{},
				},
				{
					url:      "https://subdomain.crawler-test.com",
					children: []*Node{},
				},
				{
					url:      "https://invalid.crawler-test.com",
					children: []*Node{},
				},
				{
					url:      "http://crawler-test.com/",
					children: []*Node{},
				},
				{
					url:      "http://crawler-test.com",
					children: []*Node{},
				},
				{
					url:      "https://crawler-test.com",
					children: []*Node{},
				},
			},
			err: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := NewTree(testUrl)
			doc, err := get(root.Root.url)

			tt.err(t, err)
			assert.Equalf(t, tt.want, getLinks(doc), "getLinks(%v)", doc)
		})
	}
}
