package utils

import (
	"fmt"
	"os"
	"time"
)

type Dump struct {
	Language string
	Kind     string
	Tense    string
	Wpm      int
}

func Export(data Dump) (exported bool) {

	exportPath := "./data/conju.csv"

	now := time.Now().Format(time.RFC3339)

	dataStr := fmt.Sprintf(
		"%s,%s,%s,%s,%d",
		now,
		data.Language,
		data.Kind,
		data.Tense,
		data.Wpm,
	)

	var updatedData string

	currData, err := os.ReadFile(exportPath)

	if err == nil {
		updatedData = fmt.Sprintf("%s\n%s", string(currData), dataStr)
	} else {
		headerStr := "time,language,tense,word_per_minute"
		updatedData = fmt.Sprintf(
			"%s\n%s",
			headerStr,
			dataStr,
		)
	}

	f, err := os.Create(exportPath)
	if err != nil {
		panic("Could Not Create Conju.csv")
	}

	defer f.Close()

	_, err = f.WriteString(updatedData)
	if err != nil {
		panic("Could Not Write to Conju.csv")
	}

	f.Sync()

	return true
}
