package main

import (
	"flag"
	"log"
	"os"

	"strings"

	"github.com/natefinch/lumberjack"
)

const defaultFileName = "./tmpfile.html"

var flagURL = flag.String("url", "", "(required) url to fetch")
var flagLogfile = flag.String("log", "", "(optional) error log file name")

func main() {
	flag.Parse()

	logfilename := *flagLogfile
	url := *flagURL

	if "" == url {
		flag.Usage()
		return
	}

	if "" != logfilename {
		log.SetOutput(&lumberjack.Logger{
			Filename:   logfilename,
			MaxSize:    1, // megabytes
			MaxBackups: 3,
			MaxAge:     7,     //days
			Compress:   false, // disabled by default
		})
	}

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

	delimeters := []rune{
		' ', ',', '.', '!', '?', '"', ';', ':',
		'[', ']', '(', ')', '\n', '\r', '\t', '«',
		'»', '—', '–', '“', '”', '…', '°', '²', '³',
	}

	words := Words{}
	words.Load(&buf, delimeters)

	uniq := words.GetUniq()

	err = words.PrintUniq(os.Stdout, uniq)
	if err != nil {
		log.Fatalln(err)
	}
}
