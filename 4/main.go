package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	f, err := os.Open("ex4.html")
	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()

	nodes, err := html.Parse(f)
	if err != nil {
		fmt.Println(err)
	}

	links := []link{}
	for n := range nodes.Descendants() {
		if !isLink(n.Data) {
			continue
		}

		href, ok := getHref(n)
		if !ok {
			continue
		}

		text := getLinkText(n)
		links = append(links, link{
			href: href,
			text: text,
		})
	}
	fmt.Printf("%+v\n", links)
}

type link struct {
	href string
	text string
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

	return text.String()
}
