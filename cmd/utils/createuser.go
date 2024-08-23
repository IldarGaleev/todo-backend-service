package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/IldarGaleev/todo-backend-service/internal/storage/postgresdb"
	"golang.org/x/crypto/bcrypt"

	configApp "github.com/IldarGaleev/todo-backend-service/internal/app/configapp"
)

type createUserOptions struct {
	Username string
	Password string
}

func getOptions() createUserOptions {

	username := flag.String("username", "", "enter username")
	password := flag.String("password", "", "user password")

	flag.Parse()

	return createUserOptions{
		Username: *username,
		Password: *password,
	}
}

func createTextLogger() *slog.Logger {
	log := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)

	return log
}

func main() {

	log := createTextLogger()
	options := getOptions()

	if options.Username == "" {
		for {
			fmt.Print("username: ")
			n, err := fmt.Scanln(&options.Username)

			if n == 0 {
				err = fmt.Errorf("empty username")
			}

			if err != nil {
				fmt.Println(err)
				fmt.Println("Try again")
				continue
			}
			break
		}
	}

	if options.Password == "" {
		for {
			fmt.Print("password: ")
			n, err := fmt.Scanln(&options.Password)

			if n == 0 {
				err = fmt.Errorf("empty password")
			}

			if err != nil {
				fmt.Println(err)
				fmt.Println("Try again")
				continue
			}
			break
		}
	}

	confPath := "config.yml"
	//Init app config
	appConf := configApp.MustLoadConfig(confPath)

	storageProvider := postgresdb.New(log, appConf.Dsn)
	storageProvider.MustRun()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	var wg sync.WaitGroup

	wg.Add(1)
	go func(opt createUserOptions) {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(opt.Password), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		newUser, err := storageProvider.CreateAccount(ctx, opt.Username, passwordHash)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Created new user with id: %d", *newUser.UserID)
		wg.Done()
	}(options)

	wg.Wait()
	storageProvider.Stop()
	cancel()

}
