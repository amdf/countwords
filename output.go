package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type wordCount struct {
	Word  string
	Count int
}

//PrintUniq PrintUniq
func (w Words) PrintUniq(uniq map[string]int) {
	db, err := sql.Open("sqlite", ":memory:")

	if err != nil {
		log.Fatalln("db fail")
	} else {
		log.Println("db ok")
	}

	str := `CREATE TABLE words
	(
		word TEXT NOT NULL,
		count INTEGER NOT NULL
	)`
	_, err = db.Exec(str)
	if err != nil {
		log.Fatalln(err)
	}

	str = `INSERT INTO words (word, count) VALUES ($1, $2)`
	for k, v := range uniq {
		_, err = db.Exec(str, k, v)
		if err != nil {
			log.Fatalln(err)
		}
	}

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

	for _, elem := range wrdc {
		fmt.Printf("%30s - %d\n", elem.Word, elem.Count)
	}

	db.Close()
}
