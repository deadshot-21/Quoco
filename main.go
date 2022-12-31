package main

import (
	"sync"

	"github.com/go-telegram-bot-api/api"
	"github.com/go-telegram-bot-api/bot"
	"github.com/joho/godotenv"
)

var wg sync.WaitGroup

func main() {
	godotenv.Load(".env")
	wg.Add(2)
	go bot.InitialiseBot()
	go api.InitialiseApi()
	wg.Wait()
}
