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

//type Database struct {
//    filepath    string
//    filename    string
//
//}


func QueryData(lang string, tense string) []map[string]string {

    filePath := "./data"
    filename := "spanish.db"

	connStr := fmt.Sprintf("%s/%s", filePath, filename)
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

        if currTense != "pronouns" {
            tenses = append(tenses, currTense)
        }
	}
	if err := rows.Err(); err != nil {
		panic("Could not Query Tenses in DB")
	}

	return
}

func GetPronouns(filename, filePath string) (map[string][]string) {

	connStr := fmt.Sprintf("%s/%s", filePath, filename)
	db, err := sql.Open("sqlite3", connStr)
	if err != nil {
		panic("Could Not Open DB to get Pronouns")
	}

	defer db.Close()

	if err = db.Ping(); err != nil {
        panic("Could not ping db")
	}

	rows, err := db.Query("SELECT pronoun, conjugation FROM pronouns")
	if err != nil {
		panic("Could not Query Pronouns from DB 1")
	}

	defer rows.Close()

    var pronouns = make(map[string][]string)

    type PronounSet struct {
        Pronoun string
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
