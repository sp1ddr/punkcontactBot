package punkbot

import (
	"log"
	"strconv"

	"gitlab.com/toby3d/telegraph"
)

// Account :
var Account *telegraph.Account

// err :
var err error

// TgraphToken :
var TgraphToken string

// CreateAccount : creates telegraph account
func CreateAccount() {

	Account, err = telegraph.CreateAccount(telegraph.Account{
		ShortName:  "PunkerBot",
		AuthorName: "PunkerBot",
		AuthorURL:  "https://t.me/cyberpunkrs",
	})

	log.Println("AccessToken:", Account.AccessToken)
	log.Println("AuthURL:", Account.AuthorURL)
	log.Println("ShortName:", Account.ShortName)
	log.Println("AuthorName:", Account.AuthorName)
}

// SetAccount : set object Account
func SetAccount() {

	Account = &telegraph.Account{
		AccessToken: TgraphToken,
		ShortName:   "PunkerBot",
		AuthorName:  "PunkerBot",
		AuthorURL:   "https://t.me/cyberpunkrs",
	}

}

// CreatePage :
func CreatePage(from int, name string, feedback string) (string, error) {

	content, err := telegraph.ContentFormat(feedback)
	if err != nil {
		return "", err
	}

	page, err := Account.CreatePage(telegraph.Page{
		Title:      "Report from #id" + strconv.Itoa(from) + " - " + name,
		AuthorName: Account.AuthorName,
		Content:    content,
	}, true)

	if err != nil {
		return "", err
	}

	return page.URL, nil
}
