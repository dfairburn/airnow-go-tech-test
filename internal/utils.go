package internal

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
)

func get(ctx context.Context, target string) (*goquery.Document, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("expected status code 200, got=%d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	return doc, err
}

func getLinks(baseUrl string, doc *goquery.Document) []*Node {
	var children []*Node
	if doc == nil {
		return nil
	}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		rawUrl, ok := s.Attr("href")
		if ok {
			url, err := validUrl(baseUrl, rawUrl)
			if err != nil {
				return
			}

			n := newNode(url)
			children = append(children, n)
		}
	})

	return children
}

func validUrl(base, raw string) (string, error) {
	var u *url.URL

	u, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("cannot parse url: %s", err)
	}

	if u.Host == "" {
		return validUrl(base, base+raw)
	}

	if err == nil && u.Scheme != "" && u.Host != "" {
		return u.String(), nil
	}

	return "", fmt.Errorf("cannot parse url: %s", err)
}
