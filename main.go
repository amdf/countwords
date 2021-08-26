package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"strings"

	"golang.org/x/net/html"
)

const tmpFileName = "./tmpfile.html"

var fetch = flag.String("f", "", "url to fetch")

func main() {
	flag.Parse()

	url := *fetch

	if strings.Contains(url, "http") {
		err := fetchURL(url, tmpFileName)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	rf, err := os.Open(tmpFileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	fileinfo, ferr := rf.Stat()
	if nil != ferr || 0 == fileinfo.Size() {
		fmt.Println("error: no data to analyze")
		return
	}

	fmt.Println("Begin parsing...")
	rd := bufio.NewReader(rf)
	parseHTML(rd)
}

func parseHTML(rd *bufio.Reader) {
	tokenizer := html.NewTokenizer(rd)
	prevToken := tokenizer.Token()
	var prevType string

	for {
		tt := tokenizer.Next()

		err := tokenizer.Err()
		if err == io.EOF {
			break
		}

		switch tt {
		case html.ErrorToken:
			log.Fatal(err)
		case html.StartTagToken:
			prevToken = tokenizer.Token()
			prevType = prevToken.Data

		case html.TextToken:

			if prevType == "style" || prevType == "script" || prevType == "noscript" {

			} else {
				currentToken := tokenizer.Token()
				text := strings.TrimSpace(currentToken.Data)
				if len(text) > 0 {
					fmt.Println("!!!", prevType, " === ", text)
				}
			}

		}
	}
}

func fetchURL(url string, filename string) error {
	client := http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	var wf *os.File
	wf, err = os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return err
	}

	io.Copy(wf, resp.Body)

	fileinfo, _ := wf.Stat()

	fmt.Printf("got %d bytes", fileinfo.Size())
	wf.Close()
	return nil
}
