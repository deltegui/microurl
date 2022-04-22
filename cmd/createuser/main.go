package main

import (
	"bufio"
	"fmt"
	"microurl/internal"
	"microurl/internal/config"
	"microurl/internal/persistence"
	"os"
	"strings"

	"github.com/deltegui/phoenix/hash"
	"github.com/deltegui/phoenix/validator"
)

var reader *bufio.Reader = bufio.NewReader(os.Stdin)

type CLIPresenter struct{}

func (cli CLIPresenter) Present(data interface{}) {
	fmt.Println(data)
}

func (cli CLIPresenter) PresentError(data error) {
	fmt.Printf("[ERROR] %s", data)
}

func readOrPanic() string {
	data, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(data)
}

func readUser() internal.UseCaseRequest {
	fmt.Print("Enter new user name: ")
	name := readOrPanic()
	fmt.Print("Enter the password: ")
	password := readOrPanic()
	return internal.CreateUserRequest{
		Name:     name,
		Password: password,
	}
}

func main() {
	conf := config.Load()
	conn := persistence.Connect(conf)
	conn.MigrateAll()
	repo := persistence.NewGormUserRepository(conn)
	createUser := internal.NewCreateUserCase(
		validator.New(),
		repo,
		hash.BcryptPasswordHasher{})
	request := readUser()
	createUser.Exec(CLIPresenter{}, request)
}
