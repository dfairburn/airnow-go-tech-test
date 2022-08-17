package internal

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
)

func get(target string) (*goquery.Document, error) {
	res, err := http.Get(target)
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

func getLinks(doc *goquery.Document) []*Node {
	var children []*Node
	if doc == nil {
		return nil
	}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		rawUrl, ok := s.Attr("href")
		if ok {
			url, err := validUrl(rawUrl)
			if err != nil {
				return
			}

			n := newNode(url)
			children = append(children, n)
		}
	})

	return children
}

func validUrl(raw string) (string, error) {
	var u *url.URL

	u, err := url.Parse(raw)
	if err == nil && u.Scheme != "" && u.Host != "" {
		return u.String(), nil
	}

	return "", fmt.Errorf("cannot parse url")
}
