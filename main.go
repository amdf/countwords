package main

import (
	"flag"
	"log"
	"os"

	"github.com/natefinch/lumberjack"
)

const defaultFileName = "./tmpfile.html"

var flagURL = flag.String("url", "", "(required) url to fetch")
var flagLogfile = flag.String("log", "", "(optional) error log file name")
var flagFile = flag.String("file", "", "(optional) html file to open")

func main() {
	flag.Parse()

	filename := *flagFile
	logfilename := *flagLogfile
	url := *flagURL

	if "" == url && "" == filename {
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

	if "" != url {
		filename = defaultFileName
		err := page.Create(filename, url)
		if err != nil {
			log.Fatalf("unable to create %s from %s (%s)", filename, url, err.Error())
		}
	} else {
		if "" != filename {
			err := page.Open(filename)
			if err != nil {
				log.Fatalf("unable to open %s (%s)", filename, err.Error())
			}
		}
	}

	buf, err := page.GetText()
	if err != nil {
		log.Fatalf("fail to load text (%s)", err.Error())
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
