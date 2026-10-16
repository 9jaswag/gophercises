package hparse

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func Parse(r io.Reader) ([]Link, error) {
	links := []Link{}

	nodes, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	for n := range nodes.Descendants() {
		if !isLink(n.Data) {
			continue
		}

		href, ok := getHref(n)
		if !ok {
			continue
		}

		text := getLinkText(n)
		links = append(links, Link{
			Href: href,
			Text: text,
		})
	}

	return links, nil
}

func isLink(tag string) bool {
	return tag == "a"
}

func isHref(attr string) bool {
	return attr == "href"
}

func getHref(n *html.Node) (string, bool) {
	for _, attr := range n.Attr {
		if isHref(attr.Key) {
			if attr.Val == "" {
				return "", false
			}

			return attr.Val, true
		}
	}
	return "", false
}

func getLinkText(n *html.Node) string {
	var text strings.Builder

	for sn := range n.Descendants() {
		if sn.Type == html.TextNode {
			if txt := strings.TrimSpace(sn.Data); txt != "" {
				text.WriteString(sn.Data)
			}
		}
	}

	return strings.Join(strings.Fields(text.String()), " ")
}
