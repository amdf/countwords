package main

import (
	"errors"
	"io"
	"regexp"
	"strings"
)

//Words is a word list
type Words struct {
	dsreg    *regexp.Regexp //delimeters regexp
	allWords []string
}

//Load words separated by any from the list of delimeters
func (w *Words) Load(rd io.Reader, delim []rune) (err error) {
	if nil == w {
		err = errors.New("wrong param")
		return
	}

	strbuf := new(strings.Builder)
	_, err = io.Copy(strbuf, rd)
	if err != nil {
		return
	}

	w.dsreg, err = regexp.Compile("[" + regexp.QuoteMeta(string(delim)) + "]+")
	if err != nil {
		return
	}

	w.allWords = w.dsreg.Split(strbuf.String(), -1)

	return
}

//GetUniq get map of unique words as a key and their count as a value
func (w Words) GetUniq() (result map[string]int) {
	if len(w.allWords) > 0 {
		result = make(map[string]int)
		for _, word := range w.allWords {
			upcaseWord := strings.ToUpper(word)
			count, ok := result[upcaseWord]
			if ok {
				count++
			} else {
				count = 1
			}
			result[upcaseWord] = count
		}
	}
	return
}
