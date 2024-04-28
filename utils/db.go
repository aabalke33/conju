package utils

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func QueryData(lang string, tense string) []map[string]string {
	connStr := fmt.Sprintf("./data/%s.db", lang)
	db, err := sql.Open("sqlite3", connStr)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	var verbs []map[string]string

	rows, err := db.Query(fmt.Sprintf(
		`SELECT infinitive, meaning, first_single, first_plural,
         second_single, second_plural, third_single, third_plural FROM %s`, tense))
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	type VerbTemp struct {
		infinitive    string
		meaning       string
		first_single  string
		first_plural  string
		second_single string
		second_plural string
		third_single  string
		third_plural  string
	}

	for rows.Next() {
		var tempStore VerbTemp
		if err := rows.Scan(
			&tempStore.infinitive,
			&tempStore.meaning,
			&tempStore.first_single,
			&tempStore.first_plural,
			&tempStore.second_single,
			&tempStore.second_plural,
			&tempStore.third_single,
			&tempStore.third_plural,
		); err != nil {
			log.Fatal(err)
		}

		verb := map[string]string{
			"infinitive":    tempStore.infinitive,
			"meaning":       tempStore.meaning,
			"first_single":  tempStore.first_single,
			"first_plural":  tempStore.first_plural,
			"second_single": tempStore.second_single,
			"second_plural": tempStore.second_plural,
			"third_single":  tempStore.third_single,
			"third_plural":  tempStore.third_plural,
		}

		verbs = append(verbs, verb)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return verbs
}
