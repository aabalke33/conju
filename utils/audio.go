package utils

import (
	"fmt"
	"os/exec"
)

func PlayAudio(filePath string, config Config) {

	cmd := config.FFPlayLocation
	args0 := "-autoexit"
	args1 := "-nodisp"
	args2 := "-loglevel"
	args3 := "quiet"

	err := exec.Command(cmd, args0, filePath, args1, args2, args3).Start()
	if err != nil {
		switch e := err.(type) {
		case *exec.Error:
			fmt.Println("failed executing:", err)
		case *exec.ExitError:
			fmt.Println("command exit rc =", e.ExitCode())
		default:
			fmt.Println("An Error playing audio: ", e)
		}
	}
}
