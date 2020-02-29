package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	apikey := os.Getenv("CONTACTBOT")
	bot, err := tgbotapi.NewBotAPI(apikey)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("%s", bot.Self.FirstName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.FirstName, update.Message.Text)
		command := strings.Split(update.Message.Text, " ")[0]

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.ReplyToMessageID = update.Message.MessageID
		msg.ParseMode = "Markdown"
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		switch command {
		case "/feedback":
			splited := strings.Split(update.Message.Text, " ")
			feedbackstr := strings.Join(splited[1:], " ")

			// check len of feedback message
			if len(feedbackstr) >= 10 && len(feedbackstr) <= 400 {
				reportMsg := fmt.Sprintf(" { #Feedbacks }\nFrom: [%s](tg://user?id=%d)\n\n*Feedback*: %s", update.Message.From.FirstName, update.Message.From.ID, feedbackstr)

				messageFeedback := tgbotapi.NewMessage(-1001296144335, reportMsg)
				messageFeedback.ParseMode = "Markdown"

				bot.Send(messageFeedback)
			}
		case "/afk":
			msg.Text = update.Message.From.FirstName + " *está afk*"

		case "/back":
			msg.Text = update.Message.From.FirstName + " *está de volta*"
		case "/punker":
			msg.Text = "*Hello punker!*"
		case "/status":
			msg.Text = "*online*"
		case "/ban":
			var kickobj tgbotapi.KickChatMemberConfig
			kickobj.UserID = update.Message.ReplyToMessage.From.ID
			kickobj.ChatID = update.Message.Chat.ID
			bot.KickChatMember(kickobj)
			msg.Text = "Usuário " + update.Message.ReplyToMessage.From.FirstName + " foi banido(a)"

		default:
			continue
		}
		bot.Send(msg)

	}
}
