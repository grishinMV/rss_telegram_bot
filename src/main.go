package main

import (
	"fmt"
	"net/http"
	"os"
	"rss-bot/src/db"
	"rss-bot/src/events"
	"rss-bot/src/loggrer"
	"rss-bot/src/repository"
	"rss-bot/src/telegram"
	"time"
)

func main() {
	telegramClient := telegram.NewTelegramClient("https://api.telegram.org", os.Getenv("API_KEY"), http.DefaultClient)
	dbConnection, err := db.GetConnection(db.LoadConfigFromEnv())
	if err != nil {
		fmt.Println("Error " + err.Error())
		os.Exit(1)
	}

	//feedRepository := repository.NewFeedRepository(dbConnection)
	usersRepository := repository.NewUsersRepository(dbConnection)
	logger := &loggrer.Logger{}
	eventManager := events.NewEventManager(logger)
	//parser := NewParser(eventManager, feedRepository, &http.Client{})
	registerHandlers(eventManager, telegramClient, logger, usersRepository)
	var lastMessageId int

	for {
		updates, err := telegramClient.GetUpdates(lastMessageId, 100)
		if err != nil {
			continue
		}

		for _, update := range updates.Result {
			lastMessageId = update.UpdateID + 1
			go eventManager.Dispatch(events.NewMessage{
				Message: update.Message,
			})
		}

		//feeds, _ := feedRepository.FindForUpdate()
		//
		//for _, feed := range feeds {
		//	go parser.ParseFeed(feed)
		//}

		time.Sleep(2 * time.Second)
	}
}

func registerHandlers(
	em *events.Manager,
	messenger *telegram.Client,
	logger *loggrer.Logger,
	usersRepository *repository.UsersRepository,
) {
	feedUpdatedHandler := events.NewFeedUpdatedHandler(messenger, logger, usersRepository)
	em.RegisterHandler(feedUpdatedHandler)

	newMessageHandler := events.NewNewMessageHandler(logger, usersRepository, messenger)
	em.RegisterHandler(newMessageHandler)
}
