package main

import (
	"flag"
	"fmt"
	"os"
)

var commands = map[string]func(api){
	"all":    login(getAll),
	"delete": login(delete),
	"create": login(create),
	"qr":     login(generateQR),
}

func errd(format string, args ...interface{}) {
	err(1, format, args...)
}

func err(code int, format string, args ...interface{}) {
	fmt.Printf(format, args...)
	os.Exit(code)
}

var (
	service string
	command string
	name    string
	value   string
	id      uint
)

func init() {
	flag.StringVar(&service, "url", "", "Service URL.")
	flag.StringVar(&command, "cmd", "", "Command to execute: login, all, delete, create, qr.")
	flag.StringVar(&name, "n", "", "URL name")
	flag.StringVar(&value, "v", "", "URL Value")
	flag.UintVar(&id, "id", 0, "URL id")
}

func main() {
	flag.Parse()
	cmd, ok := commands[command]
	if !ok {
		err(5, "Invalid command: %s\n", command)
	}
	api := api{service, ""}
	cmd(api)
}

func login(next func(api)) func(api) {
	return func(api api) {
		username := os.Getenv("MURL_USER")
		password := os.Getenv("MURL_PASS")
		if len(username) < 3 || len(password) < 3 {
			err(2, "Username or password are not long enough (min 3 characters). Remember to use env variables MURL_USER and MURL_PASS.\n")
		}
		api.login(username, password)
		next(api)
	}
}

func getAll(api api) {
	urls := api.getAll()
	for _, u := range urls {
		u.Print()
	}
}

func delete(api api) {
	id := readID()
	api.delete(id).Print()
}

func create(api api) {
	if len(name) <= 3 || len(value) <= 3 {
		err(4, "URL's name or value are not long enough (min 3 characters)\n")
	}
	api.create(name, value).Print()
}

func generateQR(api api) {
	id := readID()
	fmt.Println(api.generateQR(id))
}

func readID() uint {
	if id == 0 {
		err(3, "Invalid URL id: %d\n", id)
	}
	return id
}
