package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/user"
)

func path() string {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s/.murl", user.HomeDir)
}

func saveToken(token string) {
	var (
		file *os.File
		err  error
	)
	if tokenDontExist() {
		if file, err = os.Create(path()); err != nil {
			errd("Cannot create token file: %s\n", err.Error())
		}
	} else {
		if file, err = os.Open(path()); err != nil {
			errd("Cannot open token file: %s\n", err.Error())
		}
	}
	defer file.Close()
	if _, err := file.WriteString(token); err != nil {
		errd("Cannot wire to token file: %s\n", err.Error())
	}
}

func tokenDontExist() bool {
	_, err := os.Stat(path())
	return errors.Is(err, os.ErrNotExist)
}

func readToken() string {
	if tokenDontExist() {
		errd("You need to login to server:")
	}
	file, err := os.Open(path())
	if err != nil {
		errd("Cannot open token file: %s\n", err.Error())
	}
	reader := bufio.NewReader(file)
	raw, _, err := reader.ReadLine()
	if err != nil {
		errd("Cannot read from token file: %s\n", err.Error())
	}
	return string(raw)
}
