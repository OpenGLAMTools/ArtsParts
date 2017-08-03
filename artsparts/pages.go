package main

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type Page struct {
	Path  string
	Title string `yaml:"title"`
	Text  string `yaml:"text"`
}

func loadPages(fpath string) ([]Page, error) {
	f, err := os.Open(fpath)
	return parsePages(f), err
}

func parsePages(r io.Reader) []Page {
	var pages []Page
	var page Page
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		l := scanner.Text()
		switch {
		case strings.HasPrefix(l, "newpage: "):
			pages = append(pages, page)
			page = Page{}
			page.Path = strings.TrimPrefix(l, "newpage: ")
		case strings.HasPrefix(l, "# "):
			page.Title = strings.TrimPrefix(l, "# ")
		default:
			page.Text = page.Text + l + "\n"
		}

	}
	pages = append(pages, page)
	return pages[1:]
}

func getPage(pagePath string, pages []Page) Page {
	for _, p := range pages {
		if pagePath == p.Path {
			return p
		}
	}
	return Page{}
}
