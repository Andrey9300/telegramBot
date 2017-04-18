package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"log"
	"net/http"
	"gopkg.in/telegram-bot-api.v4"
)

type Joke struct{
	ID uint32 `json:"id"`
	Joke string `json:"joke"`
}

type JokeResponse struct{
	Type string `json:"type"`
	Value Joke `json:"value"`
}

var buttons = []tgbotapi.KeyboardButton{
	tgbotapi.KeyboardButton{Text: "Получить шутку"},
}

const WebhookUrl = "https://lopatintelegrambot.herokuapp.com/"

func getJoke() string{
	c := http.Client{}
	resp, err := c.Get("http://api.icndb.com/jokes/random/?limitTo=[nerdy]")
	if err != nil{
		return "jokes api is not responding"
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	joke:= JokeResponse{}

	err = json.Unmarshal(body, &joke)
	if err != nil{
		return "Joke error"
	}

	return joke.Value.Joke
}

func main() {
	port := os.Getenv("PORT")
	bot, err := tgbotapi.NewBotAPI("367470504:AAG-UJ5U_7-tfhCmkqzCAm4M6U_4r3S1ycc")
	if err != nil{
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Auth on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookUrl))

	if err != nil{
		log.Fatal(err)
	}

	updates := bot.ListenForWebhook("/")

	go http.ListenAndServe(":" + port, nil)

	for update := range updates{
		var message tgbotapi.MessageConfig
		log.Println("received text: ", update.Message.Text)

		switch update.Message.Text{
		case "Получить шутку":
			message = tgbotapi.NewMessage(update.Message.Chat.ID, getJoke())
		default:
			message = tgbotapi.NewMessage(update.Message.Chat.ID, `Нажмите "Получить шутку"`)
		}
		log.Println("received text: ", buttons)
		message.ReplyMarkup = tgbotapi.NewReplyKeyboard(buttons)

		bot.Send(message)
	}
}
