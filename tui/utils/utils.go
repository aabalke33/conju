package utils

import (
	"database/sql"
	"fmt"
	"math/rand"
	"log"

	_ "github.com/mattn/go-sqlite3"
)
 
 //func main() {
 //    wpm := int(float32(count) / float32(user.Seconds) * 60)
 //    fmt.Printf("\r\033[K")
 //    fmt.Printf("Completed %d in %d seconds. %d words per minute.", count, user.Seconds, wpm)
 //}

func ChooseVerb(
    verbs []map[string]string, keepPronouns map[string]bool)(
        verb map[string]string, pov, pronoun string) {
 
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
 
    idxVerb := rand.Int() % len(verbs)
    idxPov := rand.Int() % len(povs)

    verb = verbs[idxVerb]
    pov = povs[idxPov]
    idxPronoun := rand.Int() % len(pronouns[pov])
    pronoun = pronouns[pov][idxPronoun]

    return verb, pov, pronoun
}

func QueryData(lang string, tense string) []map[string]string {
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
         log.Fatal(err)
     }
 
     return verbs
 }
