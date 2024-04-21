// Requires: Bring your own accents
package main

import (
	"bufio"
    "io"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	s "strings"
	"time"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
    Language string `json:"language"`
    Tense string `json:"tense"`
    Seconds int `json:"seconds"`
    Pronouns map[string]bool `json:"pronouns"`
}

func readUser(path string) User {
    file, err := os.Open(path)
    if err != nil {
        log.Fatal(err)
    }

    defer file.Close()

    bytes, _ := io.ReadAll(file)

    var user User
    json.Unmarshal(bytes, &user)
    return user
}

func main() {
    user := readUser("../data/user.json")
    verbs := queryData(user.Language, user.Tense)
    haveTime := make(chan bool)
    count := 0

    fmt.Printf("\n\033[F")

    go timer(haveTime, user.Seconds)
    go game(&count, verbs, user.Pronouns)
    <- haveTime

    wpm := int(float32(count) / float32(user.Seconds) * 60)
    fmt.Printf("\r\033[K")
    fmt.Printf("Completed %d in %d seconds. %d words per minute.", count, user.Seconds, wpm)
}

func timer(timer chan<- bool, seconds int) {
    duration := time.Duration(seconds) * time.Second
    time.Sleep(duration)
    timer <- true
}

func game(count *int, verbs []map[string]string, keepPronouns map[string]bool) {

    var povs []string

    for pronoun, keepPronoun := range keepPronouns {
        if keepPronoun { povs = append(povs, pronoun) }
    }

    pronouns := map[string][]string{
        "first_single":  {"yo"},
        "first_plural":  {"nosotros"},
        "second_single": {"tu"},
        "second_plural": {"vosotros"},
        "third_single":  {"él", "ella", "usted"},
        "third_plural":  {"ellos", "ellas", "ustedes"},
    }

    for {
        idxVerb := rand.Int() % len(verbs)
        idxPov := rand.Int() % len(povs)

        verb := verbs[idxVerb]
        pov := povs[idxPov]
        idxPronoun := rand.Int() % len(pronouns[pov])
        pronoun := pronouns[pov][idxPronoun]

        if round(verb, pov, pronoun) {
            *count++
        }
    }
}

func round(verb map[string]string, pov string, pronoun string) bool {
    prompt := fmt.Sprintf("%s %s: ", pronoun, verb["infinitive"])
    fmt.Printf(prompt)

    scanner := bufio.NewScanner(os.Stdin)

    for scanner.Scan() {
        
        input := scanner.Text()

        fmt.Printf("\033[F\033[K")

        switch s.ToLower(input) {
            case s.ToLower(verb[pov]): {
                playAudio(".\\resources\\pass.mp3")
                return true
            }
            case "": return false
        }

        playAudio(".\\resources\\fail.mp3")
        fmt.Printf(prompt)
    }
    return false
}

func queryData(lang string, tense string) []map[string]string {
    connStr := fmt.Sprintf("../data/%s.db", lang)
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
        infinitive string
        meaning string
        first_single string
        first_plural string
        second_single string
        second_plural string
        third_single string
        third_plural string
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

        verb := map[string]string {
            "infinitive": tempStore.infinitive,
            "meaning": tempStore.meaning,
            "first_single": tempStore.first_single,
            "first_plural": tempStore.first_plural,
            "second_single": tempStore.second_single,
            "second_plural": tempStore.second_plural,
            "third_single": tempStore.third_single,
            "third_plural": tempStore.third_plural,
        }

        verbs = append(verbs, verb)
    }
    if err := rows.Err(); err != nil {
        log.Fatal()
    }

    return verbs
}
