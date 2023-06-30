package db

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

func passwordFromCommand(cmd string) (string, error) {
	res, err := exec.Command(cmd).Output()
	r := string(res)
	return strings.TrimSpace(r), err
}

func passwordFromShell() (string, error) {
	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	r := string(bytePassword)
	return strings.TrimSpace(r), err
}
