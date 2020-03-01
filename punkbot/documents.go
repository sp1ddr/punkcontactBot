package punkbot

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// FileObject :
var FileObject tgbotapi.FileConfig

// TelegramToken :
var TelegramToken string

// JSONResponse : generic api response
var JSONResponse map[string]string

// IpfsRouter : endpoint for Node
var IpfsRouter string

// Bot : Object
var Bot *tgbotapi.BotAPI

func init() {
	IpfsRouter = "http://127.0.0.1:1984"
}

// SaveFileToDisk : from tg to disk
func SaveFileToDisk(update tgbotapi.Update) {

	ChatID := update.Message.Chat.ID
	msg := tgbotapi.NewMessage(ChatID, "")
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = "Markdown"

	fileID := update.Message.Document.FileID
	//fileName := update.Message.Document.FileName
	FileObject.FileID = fileID

	// Dl file from Telegram Server and Save File to Disk
	file, _ := Bot.GetFile(FileObject)
	uri := file.Link(TelegramToken)
	filename := strings.Split(file.FilePath, "/")[1]
	filename = "./storage/" + filename

	if err := DownloadFile(filename, uri); err != nil {
		msg.Text = "Err downloading file from telegram server"
		Bot.Send(msg)
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
		Bot.Send(msg)
	} else {
		body := &bytes.Buffer{}
		_, _ = body.ReadFrom(response.Body)

		_ = json.NewDecoder(body).Decode(&JSONResponse)
		msg.Text = "*IPFS Hash:* `" + JSONResponse["hash"] + "`"
		Bot.Send(msg)
	}
}
