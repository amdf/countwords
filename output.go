package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

type wordCount struct {
	Word  string
	Count int
}

//PrintUniq prints unique words and their count to file, using temp DB for sorting
func (w Words) PrintUniq(f *os.File, uniq map[string]int) (err error) {
	var db *sql.DB
	db, err = w.getTemporaryDB()

	if err != nil {
		err = errors.New("fail to init temp DB")
		return
	}
	defer db.Close()

	//insert words
	str := `INSERT INTO words (word, count) VALUES ($1, $2)`
	for k, v := range uniq {
		_, err = db.Exec(str, k, v)
		if err != nil {
			log.Fatalln(err)
		}
	}

	//sort by word count
	var wrdc []wordCount

	rows, dberr := db.Query(`SELECT word, count FROM words ORDER BY count ASC`)

	if dberr == nil {
		wr := wordCount{}
		for rows.Next() {
			dberr = rows.Scan(&wr.Word, &wr.Count)
			wrdc = append(wrdc, wr)
		}
		rows.Close()
	}

	//print sorted
	for _, elem := range wrdc {
		fmt.Fprintf(f, "%30s - %d\n", elem.Word, elem.Count)
	}

	return
}

func (w Words) getTemporaryDB() (db *sql.DB, err error) {
	db, err = sql.Open("sqlite", ":memory:")

	if err != nil {
		return
	}

	str := `CREATE TABLE words
	(
		word TEXT NOT NULL,
		count INTEGER NOT NULL
	)`
	_, err = db.Exec(str)

	return
}
