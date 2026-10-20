package sitemap

import (
	"bytes"
	"encoding/xml"
	"fmt"
	hparse "hparse/parse"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type Url struct {
	Loc string `xml:"loc"`
}
type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	URLs    []Url    `xml:"url"`
}

func WriteSitemap(w io.Writer, url string) error {
	log.Println("=== generating sitemap for", url)
	siteLinks, err := getSiteLinks(url)
	if err != nil {
		fmt.Println(err)
	}

	sitemap := generateSitemap(siteLinks)
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")

	if _, err := io.WriteString(w, xml.Header); err != nil {
		return err
	}

	if err := enc.Encode(&sitemap); err != nil {
		return err
	}

	_, err = io.WriteString(w, "\n")
	return err
}

func getSiteLinks(url string) ([]string, error) {
	seenLinks := make(map[string]struct{})

	var crawl func(link string) ([]string, error)
	crawl = func(link string) ([]string, error) {
		if _, ok := seenLinks[link]; ok {
			return nil, nil
		}

		seenLinks[link] = struct{}{}

		pageLinks, err := getPageLinks(link)
		if err != nil {
			return nil, err
		}

		links := []string{link}

		for _, nextLink := range pageLinks {
			childLinks, err := crawl(nextLink)
			if err != nil {
				log.Println("Error occurred while crawling", nextLink)
				continue
			}

			links = append(links, childLinks...)
		}

		return links, nil
	}

	links, err := crawl(url)
	if err != nil {
		return nil, err
	}

	return links, nil
}

func getPageLinks(url string) ([]string, error) {
	mkSlice, err := getPageMarkup(url)
	if err != nil {
		return nil, err
	}

	links, err := getMarkupLinks(mkSlice, url)
	if err != nil {
		return nil, err
	}

	return links, nil
}

func getPageMarkup(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func getMarkupLinks(mkSlice []byte, url string) ([]string, error) {
	br := bytes.NewReader(mkSlice)
	links, err := hparse.Parse(br)
	if err != nil {
		return nil, err
	}
	ls := []string{}
	seen := make(map[string]struct{})

	for _, l := range links {
		href, err := normaliseLink(url, l.Href)
		if err != nil {
			// return nil, err
			// should we return this links err?
			continue
		}

		if _, exists := seen[href]; exists {
			continue
		}

		if sameDomain(url, href) {
			seen[href] = struct{}{}
			if ok := isDocument(href); !ok {
				ls = append(ls, href)
			}
		}
	}

	return ls, nil
}

func normaliseLink(b, l string) (string, error) {
	base, err := url.Parse(b)
	if err != nil {
		return "", err
	}

	lu, err := url.Parse(l)
	if err != nil {
		return "", err
	}

	nl := base.ResolveReference(lu)

	if nl.Scheme == "http" || nl.Scheme == "https" {
		nl.Scheme = "https"
	}
	if nl.Path != "/" {
		nl.Path = strings.TrimRight(nl.Path, "/")
		nl.RawPath = ""
	}
	nl.Host = strings.ToLower(nl.Host)

	nl.Fragment = ""
	nl.RawFragment = ""

	return nl.String(), nil
}

func sameDomain(a, b string) bool {
	u1, err1 := url.Parse(a)
	u2, err2 := url.Parse(b)

	if err1 != nil || err2 != nil {
		return false
	}

	return u1.Hostname() == u2.Hostname()
}

func isDocument(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	switch strings.ToLower(path.Ext(u.Path)) {
	case ".pdf", ".doc", ".docx", ".xls", ".xlsx",
		".ppt", ".pptx", ".odt", ".ods", ".txt":
		return true
	default:
		return false
	}
}

func generateSitemap(links []string) URLSet {
	sitemap := URLSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
	}

	for _, l := range links {
		sitemap.URLs = append(sitemap.URLs, Url{
			Loc: l,
		})
	}

	return sitemap
}
