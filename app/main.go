// Requires: Bring your own accents
package main

import (
    "os"
	"bufio"
	"database/sql"
	"fmt"
	"log"
    s "strings"
	"math/rand"
	"time"
	_ "github.com/mattn/go-sqlite3"
)

type Inputs struct {
    lang string
    tense string
    seconds int
    keepVosotros bool
}

func main() {
    var inputs = Inputs{
        "spanish",
        "present",
        60,
        true,
    }

    verbs := queryData(inputs.lang, inputs.tense)

    haveTime := make(chan bool)
    duration := time.Duration(inputs.seconds) * time.Second

    count := 0
    fmt.Printf("\n\033[F")

    go timer(haveTime, duration)
    go game(&count, verbs)

    <- haveTime

    wpm := int(count * inputs.seconds / 60)
    fmt.Printf("\r\033[K")
    fmt.Printf("Completed %d in %d seconds. %d words per minute.", count, inputs.seconds, wpm)
}

func timer(timer chan<- bool, duration time.Duration) {
    time.Sleep(duration)
    timer <- true
}

func game(count *int, verbs []map[string]string) {

    povs := [6]string {
        "first_single",
        "first_plural",
        "second_single",
        "second_plural",
        "third_single",
        "third_plural",
    }

    pronouns := map[string][3]string{
        "first_single": {"yo"},
        "first_plural": {"nosotros"},
        "second_single": {"tu"},
        "second_plural": {"vosotros"},
        "third_single": {"él", "ella", "usted"},
        "third_plural": {"ellos", "ellas", "ustedes"},
    }

    for {
        idxVerb := rand.Int() % len(verbs)
        idxTense := rand.Int() % len(povs)
        verb := verbs[idxVerb]
        pov := povs[idxTense]
        pronoun := getPronoun(pronouns, pov)

        if round(verb, pov, pronoun) {
            *count++
        }
    }
}

func getPronoun(pronouns map[string][3]string, pov string) string {
    if pov == "third_single" || pov == "third_plural" {
        idx := rand.Int() % 3
        return pronouns[pov][idx]
    }
    return pronouns[pov][0]
}

func round(verb map[string]string, pov string, pronoun string) bool {
    prompt := fmt.Sprintf("%s %s: ", pronoun, verb["infinitive"])
    fmt.Printf(prompt)

    scanner := bufio.NewScanner(os.Stdin)

    for scanner.Scan() {
        
        input := scanner.Text()

        fmt.Printf("\033[F\033[K")

        switch s.ToLower(input) {
            case s.ToLower(verb[pov]): return true
            case "": return false
        }

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

    rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", tense))
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
            &tempStore.second_single,
            &tempStore.third_single,
            &tempStore.first_plural,
            &tempStore.second_plural,
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
