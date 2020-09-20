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
)

var (
	token = ""
)

func init() {
	if err := godotenv.Load(fmt.Sprintf("./%s.env", os.Getenv("GO_ENV"))); err != nil {
		log.Fatal(err)
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
