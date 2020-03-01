package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sp1ddr/punkcontactBot/punkbot"
)

// IpfsRouter : endpoint
var IpfsRouter string

// JSONResponse : generic api response
var JSONResponse map[string]string

func main() {

	IpfsRouter = "http://127.0.0.1:1984"

	tgapikey := os.Getenv("CONTACTBOT")
	tgraphapikey := os.Getenv("TELEGRAPH")

	if tgapikey == "" {
		fmt.Println("run 'export CONTACTBOT=your_api_key' ")
		os.Exit(1)
	}

	if tgraphapikey == "" {
		fmt.Println("[ warning ] telegraph api key not informed! The bot will create one right now.")
		punkbot.CreateAccount()

		fmt.Println("\t Telegraph key:  ", punkbot.Account.AccessToken, " run 'export TELEGRAPH=your_api_key' for full integration")

		punkbot.TgraphToken = punkbot.Account.AccessToken
	}

	punkbot.TgraphToken = tgraphapikey
	punkbot.SetAccount()

	bot, err := tgbotapi.NewBotAPI(tgapikey)
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

		var FileObject tgbotapi.FileConfig

		ChatID := update.Message.Chat.ID

		msg := tgbotapi.NewMessage(ChatID, "")
		msg.ReplyToMessageID = update.Message.MessageID
		msg.ParseMode = "Markdown"
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message.Document != nil {
			fileID := update.Message.Document.FileID
			//fileName := update.Message.Document.FileName
			FileObject.FileID = fileID

			// Dl file from Telegram Server and Save File to Disk
			file, _ := bot.GetFile(FileObject)
			uri := file.Link(tgapikey)
			filename := strings.Split(file.FilePath, "/")[1]
			filename = "./storage/" + filename

			if err := punkbot.DownloadFile(filename, uri); err != nil {
				msg.Text = "Err downloading file from telegram server"
				bot.Send(msg)
			}

			f, _ := os.Open(filename)
			defer f.Close()

			// Send File to IPFS-Router
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, _ := writer.CreateFormFile("file", filepath.Base(f.Name()))

			io.Copy(part, f)
			writer.Close()

			req, err := http.NewRequest("POST", IpfsRouter+"/ipfs/file", body)
			req.Header.Add("Content-Type", writer.FormDataContentType())

			client := &http.Client{}
			response, err := client.Do(req)
			if err != nil {
				msg.Text = "Err sending to ipfs-router"
				bot.Send(msg)
			} else {
				body := &bytes.Buffer{}
				_, _ = body.ReadFrom(response.Body)

				_ = json.NewDecoder(body).Decode(&JSONResponse)
				msg.Text = "*IPFS Hash:* `" + JSONResponse["hash"] + "`"
				bot.Send(msg)
			}

		}

		switch command {
		case "/feedback":
			if update.Message.Chat.Type != "group" && update.Message.Chat.Type != "supergroup" {
				splited := strings.Split(update.Message.Text, " ")
				feedbackstr := strings.Join(splited[1:], " ")

				// check len of feedback message
				if len(feedbackstr) >= 10 && len(feedbackstr) <= 400 {

					uri, err := punkbot.CreatePage(update.Message.From.ID, update.Message.From.FirstName, feedbackstr)
					if err != nil {
						msg.Text = "*Error savin data on telegraph*"
						bot.Send(msg)
					}

					reportMsg := fmt.Sprintf(" { #Feedbacks }\nFrom: [%s](tg://user?id=%d) `#id%d`\n\n*Feedback*: %s", update.Message.From.FirstName, update.Message.From.ID, update.Message.From.ID, uri)

					messageFeedback := tgbotapi.NewMessage(-1001296144335, reportMsg)
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
		default:
			continue
		}

	}
}
