package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gomarkdown/markdown"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sp1ddr/punkcontactBot/punkbot"
)

// IpfsRouter : endpoint
var IpfsRouter string

// JSONResponse : generic api response
var JSONResponse map[string]string

func main() {

	tgapikey := os.Getenv("CONTACTBOT")
	tgraphapikey := os.Getenv("TELEGRAPH")

	if tgapikey == "" {
		fmt.Println("run 'export CONTACTBOT=your_api_key' ")
		os.Exit(1)
	}

	punkbot.TelegramToken = tgapikey

	if tgraphapikey == "" {
		fmt.Println("[ warning ] telegraph api key not informed! The bot will create one right now.")
		punkbot.CreateAccount()

		fmt.Println("\t Telegraph key:  ", punkbot.Account.AccessToken, " run 'export TELEGRAPH=your_api_key' for full integration")

		punkbot.TgraphToken = punkbot.Account.AccessToken
	} else {
		punkbot.TgraphToken = tgraphapikey
	}

	punkbot.SetAccount()

	bot, err := tgbotapi.NewBotAPI(tgapikey)
	if err != nil {
		log.Panic(err)
	}

	// Set bot pointer to other references for manipulation
	punkbot.Bot = bot

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

		ChatID := update.Message.Chat.ID

		msg := tgbotapi.NewMessage(ChatID, "")
		msg.ReplyToMessageID = update.Message.MessageID
		msg.ParseMode = "Markdown"
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message.Document != nil {
			punkbot.SaveFileToDisk(update)
		}

		switch command {
		case "/feedback":
			if update.Message.Chat.Type != "group" && update.Message.Chat.Type != "supergroup" {
				splited := strings.Split(update.Message.Text, " ")
				feedbackstr := strings.Join(splited[1:], " ")

				// parse feedback in markdown to html
				md := []byte(feedbackstr)
				html := markdown.ToHTML(md, nil, nil)

				// check len of feedback message
				if len(feedbackstr) >= 10 && len(feedbackstr) <= 400 {

					uri, err := punkbot.CreatePage(update.Message.From.ID, update.Message.From.FirstName, string(html))
					if err != nil {
						msg.Text = "*Error savin data on telegraph*"
						bot.Send(msg)
					}

					reportMsg := fmt.Sprintf(" { #Feedbacks }\nFrom: [%s](tg://user?id=%d) `#id%d`\n\n*Feedback*: %s", update.Message.From.FirstName, update.Message.From.ID, update.Message.From.ID, uri)

					messageFeedback := tgbotapi.NewMessage(-1001296144335, reportMsg)
					messageFeedback.ParseMode = "Markdown"
					bot.Send(messageFeedback)

					// Send to -280353697 - The Realm for debug and testing
					messageFeedback = tgbotapi.NewMessage(-280353697, reportMsg)
					messageFeedback.ParseMode = "Markdown"
					bot.Send(messageFeedback)

					msg.Text = "*Feedback enviado com sucesso!*"
					bot.Send(msg)
				} else {
					msg.Text = "_Feedback len must be greater than 10char and less than 400!_ \n*[ Error ]*"
					bot.Send(msg)
				}
			}
		case "/afk":
			msg.Text = update.Message.From.FirstName + " *está afk*"

			bot.Send(msg)
		case "/back":
			msg.Text = update.Message.From.FirstName + " *está de volta*"

			bot.Send(msg)
		case "/busy":
			msg.Text = update.Message.From.FirstName + " *está ocupado*"

			bot.Send(msg)
		case "/punker":
			msg.Text = "*Hello punker!*"

			bot.Send(msg)
		case "/status":
			msg.Text = "*online*"

			bot.Send(msg)
		case "/ban":
			var kickobj tgbotapi.KickChatMemberConfig
			kickobj.UserID = update.Message.ReplyToMessage.From.ID
			kickobj.ChatID = update.Message.Chat.ID

			bot.KickChatMember(kickobj)

			msg.Text = "Usuário " + update.Message.ReplyToMessage.From.FirstName + " foi banido(a)"
			bot.Send(msg)
		case "/pin":
			var pinobj tgbotapi.PinChatMessageConfig

			pinobj.ChatID = update.Message.Chat.ID
			pinobj.MessageID = update.Message.ReplyToMessage.MessageID
			pinobj.DisableNotification = false

			resp, _ := bot.PinChatMessage(pinobj)
			if !resp.Ok {
				msg.Text = "Err ao fixar mensagem"
			} else {
				msg.Text = "*Mensagem fixada!*"
			}
			bot.Send(msg)
		default:
			continue
		}

	}
}
