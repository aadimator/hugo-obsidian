package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// parse single file for links
func parse(dir, pathPrefix string) []Link {
	// read file
	source, err := os.ReadFile(dir)
	if err != nil {
		panic(err)
	}

	// parse md
	var links []Link
	fmt.Printf("[Parsing note] %s => ", trim(dir, pathPrefix, ".md"))

	var buf bytes.Buffer
	if err := md.Convert(source, &buf); err != nil {
		panic(err)
	}

	doc, _ := goquery.NewDocumentFromReader(&buf)
	var n int
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		target, ok := s.Attr("href")
		if !ok {
			target = "#"
		}

		target = processTarget(target)
		source := processSource(trim(dir, pathPrefix, ".md"))

		// fmt.Printf("  '%s' => %s\n", source, target)
		links = append(links, Link{
			Source: source,
			Target: target,
			Text:   text,
		})
		n++
	})
	fmt.Printf("found: %d links\n", n)

	return links
}
