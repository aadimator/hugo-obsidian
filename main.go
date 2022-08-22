package main

import (
	"flag"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	wikilink "github.com/abhinav/goldmark-wikilink"
	"github.com/yuin/goldmark"
)

var md goldmark.Markdown

func init() {
	md = goldmark.New(
		goldmark.WithExtensions(&wikilink.Extender{}),
	)
}

type Link struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Text   string `json:"text"`
}

type LinkTable = map[string][]Link
type Index struct {
	Links     LinkTable `json:"links"`
	Backlinks LinkTable `json:"backlinks"`
}

type Content struct {
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	LastModified time.Time `json:"lastmodified"`
	Tags         []string  `json:"tags"`
}

type ContentIndex = map[string]Content

type ConfigTOML struct {
	IgnoredFiles []string `toml:"ignoreFiles"`
}

func getIgnoredFiles(base string) (res map[string]struct{}) {
	res = make(map[string]struct{})

	source, err := os.ReadFile(filepath.FromSlash(base + "/config.toml"))
	if err != nil {
		return res
	}

	var config ConfigTOML
	if _, err := toml.Decode(string(source), &config); err != nil {
		return res
	}

	for _, glb := range config.IgnoredFiles {
		matches, _ := filepath.Glob(base + glb)
		for _, match := range matches {
			res[match] = struct{}{}
		}
	}

	return res
}

func main() {
	in := flag.String("input", ".", "Input Directory")
	out := flag.String("output", ".", "Output Directory")
	root := flag.String("root", "..", "Root Directory (for config parsing)")
	index := flag.Bool("index", false, "Whether to index the content")
	strip := flag.Bool("strip", true, "Whether to strip comments")
	flag.Parse()

	process_index(*in, "notes")
	ignoreBlobs := getIgnoredFiles(*root)
	l, i, mapping := walk(*in, ".md", *index, ignoreBlobs, *strip)
	f := filter(l)
	tl, ti := transform_links(f, i, mapping)
	err := write(tl, ti, *index, *out, *root)
	if err != nil {
		panic(err)
	}
}
