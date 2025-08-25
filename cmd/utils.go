package cmd

import (
	"fmt"
)

func readValue(prompt string) (string, error) {
	var dest string
	fmt.Printf("%s: ", prompt)
	_, err := fmt.Scanln(&dest)
	return dest, err
}
