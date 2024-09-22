package main

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

func setPasscode() error {
	passcode, err := os.ReadFile(".passcode")
	if err != nil {
		switch {
		case errors.Is(err, os.ErrNotExist):
			r := bufio.NewReader(os.Stdin)
			print("please enter a passcode: ")
			passcode, err := r.ReadString('\n')
			if err != nil {
				return err
			}
			passcode = strings.Trim(passcode, "\n")
			println(passcode + " written to .passcode")
			err = os.WriteFile(".passcode", []byte(passcode), 0755)
			if err != nil {
				return err
			}
			PASSCODE = passcode

			return nil

		default:
			return err
		}
	}

	PASSCODE = strings.Trim(string(passcode), "\n")

	return nil
}
