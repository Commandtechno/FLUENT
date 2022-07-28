package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ericchiang/css"
	"golang.org/x/net/html"
)

func getEmojipediaEmojis() {
	res, err := http.Get("https://emojipedia.org/microsoft-teams/")
	if err != nil {
		panic(err)
	}

	html, err := html.Parse(res.Body)
	if err != nil {
		panic(err)
	}

	sel, err := css.Parse("ul.emoji-grid > li img")
	if err != nil {
		panic(err)
	}

	for _, img := range sel.Select(html) {
		for _, attr := range img.Attr {
			if attr.Key == "data-src" {
				filename := attr.Val[strings.LastIndex(attr.Val, "/")+1:]
				fmt.Println("Downloading emoji", filename)

				file, err := os.Create(filepath.Join("emojipedia", filename))
				if err != nil {
					panic(err)
				}

				res, err := http.Get(attr.Val)
				if err != nil {
					panic(err)
				}

				if _, err := io.Copy(file, res.Body); err != nil {
					panic(err)
				}
			}
		}
	}
}