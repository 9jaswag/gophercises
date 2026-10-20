package main

import (
	"os"
	"sitemap/sitemap"
)

func main() {
	link := "https://www.calhoun.io/"

	sitemap.WriteSitemap(os.Stdout, link)
}
