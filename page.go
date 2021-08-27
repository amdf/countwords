package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

//Page is a web page
type Page struct {
	filename string
	rd       *bufio.Reader
}

//Create page. Load HTML page content from URL and open it.
func (p *Page) Create(filename string, url string) (err error) {
	if nil == p {
		err = errors.New("wrong param")
		return
	}
	client := http.Client{}

	var resp *http.Response
	resp, err = client.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var wf *os.File
	wf, err = os.Create(filename)
	if err != nil {
		return
	}
	defer wf.Close()

	var size int64
	size, err = io.Copy(wf, resp.Body)

	if size <= 0 {
		err = errors.New("no data")
	}

	err = p.Open(filename)

	return
}

//Open page. Loads HTML page from file.
func (p *Page) Open(filename string) (err error) {
	if nil == p {
		err = errors.New("wrong param")
		return
	}

	var rf *os.File
	rf, err = os.Open(filename)
	if err != nil {
		return
	}

	var fileinfo fs.FileInfo
	fileinfo, err = rf.Stat()
	if err != nil {
		return
	}

	if 0 == fileinfo.Size() {
		err = errors.New("no data")
	}

	p.rd = bufio.NewReader(rf)
	return
}

//GetText return text of the web page
func (p Page) GetText() (buf bytes.Buffer, err error) {
	if nil == p.rd {
		err = errors.New("page not loaded")
		return
	}

	tokenizer := html.NewTokenizer(p.rd)
	prevToken := tokenizer.Token()
	var prevType string

	stop := false
	for !stop {
		tt := tokenizer.Next()
		errtok := tokenizer.Err()
		if errtok != nil {
			if errtok != io.EOF {
				err = errtok
			}
			stop = true
		} else {
			switch tt {
			case html.ErrorToken:
				stop = true
			case html.StartTagToken:
				prevToken = tokenizer.Token()
				prevType = prevToken.Data

			case html.TextToken:
				if !(prevType == "style" || prevType == "script" || prevType == "noscript") {
					currentToken := tokenizer.Token()
					text := strings.TrimSpace(currentToken.Data)
					if len(text) > 0 {
						buf.WriteString(text + "\n")
					}
				}
			}
		}
	}

	return
}
