package audiodev

import (
	"bytes"
	"conju/utils"
	"encoding/json"
	"fmt"
	"math/rand"
	"os/exec"
	"strconv"
    "time"
    "bufio"
    "os"
    s "strings"
)

var jsonPath = "./audiodev/pull.json"

type Word struct {
	Text         string   `json:"text"`
	Translations []string `json:"translations"`
	AudioURL     string   `json:"audioURL"`
}

 func Test() {
     length := getLength()
     seconds := 60



     haveTime := make(chan bool)
     count := 0
 
     fmt.Printf("\n\033[F")
 
     go timer(haveTime, seconds)
     go game(&count, length)
     <- haveTime
 
     wpm := int(float32(count) / float32(seconds) * 60)
     fmt.Printf("\r\033[K")
     fmt.Printf("Completed %d in %d seconds. %d words per minute.", count, seconds, wpm)
 }
 
 func timer(timer chan<- bool, seconds int) {
     duration := time.Duration(seconds) * time.Second
     time.Sleep(duration)
     timer <- true
 }
 
 func game(count *int, length int) {
     for {
         if round(length) {
             *count++
         }
     }
 }
 
func round(length int) bool {
    word := chooseWord(length)
    config := utils.ReadConfig()
    utils.PlayAudio(word.AudioURL, config)


    var translations []string

    for i := 0; i < len(word.Translations); i++ {

        formatted := s.ReplaceAll(s.ReplaceAll(
            s.ReplaceAll(
                s.ToLower(word.Translations[i]),
                "(", ""), ")", ""), "?", "")

        translations = append(translations, formatted)

    }

    //fmt.Println(translations)

    prompt := fmt.Sprintf("%s: ", word.Text)
    fmt.Printf(prompt)

    scanner := bufio.NewScanner(os.Stdin)

    for scanner.Scan() {

        input := scanner.Text()

        fmt.Printf("\033[F\033[K")

        if input == "" {
            return false
        }

        for i := 0; i < len(translations); i++ {
            if translations[i] == input {
                //utils.playAudio(".\\resources\\pass.mp3")
                return true
            }
        }

        //playAudio(".\\resources\\fail.mp3")
        fmt.Printf(prompt)
    }
    return false
}

func getLength() (length int) {

	process := "jq"
	arg := ".words | length"
	cmd := exec.Command(process, arg, jsonPath)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	stringOutput, err := out.ReadString(0x0D)
	if err != nil {
		panic(err)
	}
	length, err = strconv.Atoi(stringOutput[:len(stringOutput)-1])
	if err != nil {
		panic(err)
	}

	return
}

func chooseWord(length int) (word Word) {

	idx := rand.Int() % length

	process := "jq"
	arg := fmt.Sprintf(".words[%d]", idx)
	cmd := exec.Command(process, arg, jsonPath)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	json.Unmarshal([]byte(out.String()), &word)

	return word
}
