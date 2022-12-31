package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-telegram-bot-api/controllers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/solywsh/chatgpt"
)

var wg sync.WaitGroup

func main() {
	godotenv.Load(".env")
	wg.Add(2)
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	// log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		chat := chatgpt.New(os.Getenv("OPENAI_KEY"), "user_id(not required)", 30*time.Second)

		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		if !update.Message.IsCommand() { // ignore any non-command Messages
			fmt.Println("MESSAGE", update.Message.Text)

			//
			//select {
			//case <-chat.GetDoneChan():
			//	fmt.Println("time out/finish")
			//}
			question := "Answer in true or false. Is the following question related to computer science? " + update.Message.Text
			log.Printf("Q: %s\n", question)
			answer, err := chat.Chat(question)
			if err != nil {
				log.Println(err)
			}
			log.Printf("A: %s\n", answer)

			if answer == "" {
				msg.Text = "I didn't understand the question!"
			} else if answer[:4] == "True" {
				c_question := update.Message.Text
				c_answer, err := chat.Chat(c_question)
				if err != nil {
					log.Println(err)
					msg.Text = "Error occured! Try again"
				} else {
					msg.Text = c_answer
				}
			} else {
				msg.Text = "Question not related to Coding!"
			}

		} else {

			switch update.Message.Command() {
			case "start":
				msg.Text = `
				Hey explorer👋 Welcome to Quoco :))
	
Need help for anything related to coding? Here you go!
I will be happy to help. Just type your question and see the magic ✨

It could take up to 30 seconds for me to answer.
Some commands other than question you can try:

/help - for more info
/status - to check the server is up or not
/joke - to get a coding joke
/quote - to get a coding quote
				`

			case "help":
				msg.Text = `
				An example for command is:

You: /quote
Quoco: "The best code is no code at all." - Unknown

An example for question is: 

You: How to sort in python?
Quoco: 

Python offers a number of ways to sort items in lists or other collections. The most common way to sort lists is with the built-in sorted() function. You can also use the list.sort() method, or you can use the built-in sorted() function.
The sorted() function takes an iterable and returns a new sorted list from that iterable. It has two optional arguments: reverse and key. The reverse argument is a boolean value; if set to True, the list is sorted in descending order. The key argument is a function that will be applied to each element before making comparisons.
For example, to sort a list of numbers in ascending order, you could use the following code:

my_list = [3, 4, 1, 5, 2]

sorted_list = sorted(my_list)  # [1, 2, 3, 4, 5]

To sort a list of strings in descending order, you could use the following code:

str_list = ['hello', 'hi', 'howdy']

sorted_str_list = sorted(str_list, reverse=True)  # ['howdy', 'hello', 'hi']
				`
			case "status":
				msg.Text = "I'm up and running."
			case "joke":
				e_answer, err := chat.Chat("tell me a coding joke")
				if err != nil {
					log.Println(err)
					msg.Text = "Error occured! Try again"
				} else {
					msg.Text = e_answer
				}
			case "quote":
				e_answer, err := chat.Chat("tell me a coding quote")
				if err != nil {
					log.Println(err)
					msg.Text = "Error occured! Try again"
				} else {
					msg.Text = e_answer
				}
			default:
				msg.Text = "Even Google can't answer everything, I'm still Quoco :("
			}
		}
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
		defer chat.Close()
	}
	r := mux.NewRouter()
	homeController := controllers.NewHomeController()
	r.HandleFunc("/", homeController.Home).Methods("GET")
	port := os.Getenv("PORT")
	fmt.Println("running on port " + port)

	http.ListenAndServe(":"+port, r)
	wg.Wait()
}
