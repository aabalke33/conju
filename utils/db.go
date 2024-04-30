package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	s "strings"

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

type Database struct {
	ProperName string
	FileName   string
}

func GetDatabases(directory string) (databases []Database) {

	c, err := os.ReadDir("./data")

	if err != nil {
		panic("Could not read data directory")
	}

	for _, entry := range c {
		isDB := s.HasSuffix(entry.Name(), ".db")
		if isDB {
			lowerName := s.ToLower(s.Replace(entry.Name(), ".db", "", -1))
			properName := s.ToUpper(string(lowerName[0])) + lowerName[1:]

			database := Database{
				ProperName: properName,
				FileName:   entry.Name(),
			}

			databases = append(databases, database)
		}
	}
	return
}

func GetTenses(filename, dataLocation string) (tenses []string) {

	connStr := fmt.Sprintf("%s/%s", dataLocation, filename)
	db, err := sql.Open("sqlite3", connStr)
	if err != nil {
		panic("Could Not Open DB to get Tenses")
	}

	defer db.Close()

	rows, err := db.Query(fmt.Sprintf(
		"SELECT name FROM sqlite_master where type='table';"))
	if err != nil {
		panic("Could not Query Tenses from DB")
	}

	defer rows.Close()

	var currTense string

	for rows.Next() {
		if err := rows.Scan(&currTense); err != nil {
			panic("Could not Query Tenses in DB")
		}

		tenses = append(tenses, currTense)
	}
	if err := rows.Err(); err != nil {
		panic("Could not Query Tenses in DB")
	}

	return
}
