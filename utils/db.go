package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	s "strings"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	ProperName        string
	LowerName         string
	Directory         string
	Basename          string
	AvailablePronouns []string
}

func (d Database) QueryData(tense string) []map[string]string {

	connStr := fmt.Sprintf("%s/%s", d.Directory, d.Basename)
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
		`SELECT * FROM %s`, tense))
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var values []interface{}

	columns, err := rows.Columns()
	if err != nil {
		panic("Could not get language columns")
	}

	for range columns {
		var i interface{}
		values = append(values, &i)
	}

	for rows.Next() {

		verb := make(map[string]string)

		if err := rows.Scan(values...); err != nil {
			log.Fatal(err)
		}

		for i, values := range values {

			deRefV := reflect.ValueOf(values).Elem().Interface()

			switch v := deRefV.(type) {
			case nil:
				verb[columns[i]] = "NULL"
			case []byte:
				verb[columns[i]] = string(v)
			case string:
				verb[columns[i]] = v
			default:
				verb[columns[i]] = fmt.Sprintf("%v", v)
			}
		}

		verbs = append(verbs, verb)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return verbs
}

func (d Database) GetTenses() (tenses []string) {
	connStr := fmt.Sprintf("%s/%s", d.Directory, d.Basename)

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

		if currTense != "pronouns" {
			tenses = append(tenses, currTense)
		}
	}
	if err := rows.Err(); err != nil {
		panic("Could not Query Tenses in DB")
	}

	return
}

func (d Database) GetPronouns(tense string, userSelectedPronouns []string) map[string][]string {

	formatStringSlice := func(slice []string) string {
		formattedStrings := make([]string, len(slice))

		for i, pronoun := range slice {
			formattedStrings[i] = fmt.Sprintf("'%s'", pronoun)
		}

		return s.Join(formattedStrings, ", ")
	}

	connStr := fmt.Sprintf("%s/%s", d.Directory, d.Basename)
	db, err := sql.Open("sqlite3", connStr)
	if err != nil {
		panic("Could Not Open DB to get Pronouns")
	}

	defer db.Close()

	if err = db.Ping(); err != nil {
		panic("Could not ping db")
	}

	query := fmt.Sprintf(`SELECT pronoun, conjugation FROM pronouns 
        WHERE conjugation IN (
	    SELECT name AS conjugation FROM pragma_table_info('%s'))
        AND conjugation IN (%s)`,
		tense, formatStringSlice(userSelectedPronouns))

	rows, err := db.Query(query)
	if err != nil {
		panic("Could not Query Pronouns from DB 1")
	}

	defer rows.Close()

	var pronouns = make(map[string][]string)

	type PronounSet struct {
		Pronoun     string
		Conjugation string
	}

	var curr PronounSet

	for rows.Next() {
		if err := rows.Scan(&curr.Pronoun, &curr.Conjugation); err != nil {
			panic("Could not Query Pronouns in DB 2")
		}

		pronouns[curr.Conjugation] = append(pronouns[curr.Conjugation], curr.Pronoun)
	}
	if err := rows.Err(); err != nil {
		panic("Could not Query Pronouns in DB 3")
	}

	return pronouns
}

func GetDatabases(directory string) map[string]Database {

	c, err := os.ReadDir(directory)
	if err != nil {
		panic("Could not read data directory")
	}

	var databases = make(map[string]Database)

	for _, entry := range c {
		isDB := s.HasSuffix(entry.Name(), ".db")
		if isDB {
			lowerName := s.ToLower(s.Replace(entry.Name(), ".db", "", -1))
			properName := s.ToUpper(string(lowerName[0])) + lowerName[1:]

			database := Database{
				ProperName: properName,
				Basename:   entry.Name(),
				LowerName:  lowerName,
				Directory:  directory,
			}

			databases[database.ProperName] = database
		}
	}
	return databases
}
