package main

import (
	"flag"
	"fmt"
	"log"

	"strings"
)

const defaultFileName = "./tmpfile.html"

var fetch = flag.String("f", "", "url to fetch")

func main() {
	flag.Parse()

	url := *fetch

	page := Page{}

	if strings.Contains(url, "http") {
		err := page.Create(defaultFileName, url)
		if err != nil {
			log.Fatalf("cannot create %s from %s (%s)", defaultFileName, url, err.Error())
		}
	} else {
		err := page.Open(defaultFileName)
		if err != nil {
			log.Fatalf("cannot open %s (%s)", defaultFileName, err.Error())
		}
	}

	buf, err := page.GetText()
	if err != nil {
		log.Fatalf("fail to load text from %s (%s)", defaultFileName, err.Error())
	}

	fmt.Print(buf.String())
	fmt.Println("\n---\n\nsize ", buf.Len())
}
