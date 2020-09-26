package main

import (
	"fmt"
	"log"
	"os"

	"github.com/daichidd/daigo-chan/cmd"

	"github.com/joho/godotenv"
)

const (
	DISCORD_TOKEN = "DISCORD_TOKEN"
	HEROKU_ENV    = "HEROKU_ENV"
)

var (
	token    = ""
	isHeroku = false
)

func init() {
	// herokuではosの環境変数しか使えない
	isHeroku = os.Getenv(HEROKU_ENV) != ""
	if !isHeroku {
		if err := godotenv.Load(fmt.Sprintf("./%s.env", os.Getenv("GO_ENV"))); err != nil {
			log.Fatal(err)
		}
	}

	token = os.Getenv(DISCORD_TOKEN)
	if token == "" {
		log.Println("token is empty")
		return
	}
}

func main() {
	if err := cmd.Run(token); err != nil {
		log.Fatal(err)
	}
}
