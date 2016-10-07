package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"gopkg.in/telegram-bot-api.v4"
)

type (
	// Repo information.
	Repo struct {
		Owner string
		Name  string
	}

	// Build information.
	Build struct {
		Event   string
		Number  int
		Commit  string
		Message string
		Branch  string
		Author  string
		Status  string
		Link    string
	}

	// Config for the plugin.
	Config struct {
		Token   string
		Debug   bool
		To      []string
		Message []string
		Photo   []string
		Format  string
	}

	// Plugin values.
	Plugin struct {
		Repo   Repo
		Build  Build
		Config Config
	}
)

func trimElement(keys []string) []string {
	var newKeys []string

	for _, value := range keys {
		value = strings.Trim(value, " ")
		if len(value) == 0 {
			continue
		}
		newKeys = append(newKeys, value)
	}

	return newKeys
}

func fileExist(keys []string) []string {
	var newKeys []string

	for _, value := range keys {
		if _, err := os.Stat(value); os.IsNotExist(err) {
			continue
		}
		newKeys = append(newKeys, value)
	}

	return newKeys
}

func parseID(keys []string) []int64 {
	var newKeys []int64

	for _, value := range keys {
		id, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			log.Println(err.Error())

			continue
		}
		newKeys = append(newKeys, id)
	}

	return newKeys
}

// Exec executes the plugin.
func (p Plugin) Exec() error {

	if len(p.Config.Token) == 0 || len(p.Config.To) == 0 {
		log.Println("missing telegram token or user list")

		return errors.New("missing telegram token or user list")
	}

	var message []string
	if len(p.Config.Message) > 0 {
		message = p.Config.Message
	} else {
		message = p.Message(p.Repo, p.Build)
	}

	bot, err := tgbotapi.NewBotAPI(p.Config.Token)

	if err != nil {
		log.Println(err.Error())

		return err
	}

	bot.Debug = p.Config.Debug

	// parse ids
	ids := parseID(p.Config.To)
	photos := fileExist(trimElement(p.Config.Photo))

	// send message.
	for _, user := range ids {
		for _, value := range trimElement(message) {
			msg := tgbotapi.NewMessage(user, value)
			msg.ParseMode = p.Config.Format

			_, err := bot.Send(msg)

			if err != nil {
				log.Println(err.Error())
			}
		}

		for _, value := range photos {
			msg := tgbotapi.NewPhotoUpload(user, value)
			_, err := bot.Send(msg)

			if err != nil {
				log.Println(err.Error())
			}
		}
	}

	return nil
}

// Message is plugin default message.
func (p Plugin) Message(repo Repo, build Build) []string {
	return []string{fmt.Sprintf("[%s] <%s> (%s)『%s』by %s",
		build.Status,
		build.Link,
		build.Branch,
		build.Message,
		build.Author,
	)}
}