package main

import (
	"fmt"
	"net/http"
	"os"
	"rss-bot/src/db"
	"rss-bot/src/events"
	"rss-bot/src/logger"
	"rss-bot/src/parser"
	"rss-bot/src/repository"
	"rss-bot/src/telegram"
	"time"
)

func main() {
	telegramClient := telegram.NewTelegramClient("https://api.telegram.org", os.Getenv("API_KEY"), http.DefaultClient)
	dbPath, exist := os.LookupEnv("DB_PATH")
	if !exist {
		dbPath = "./../rss.sqlite"
	}

	dbConnection, err := db.GetConnection(dbPath)
	if err != nil {
		fmt.Println("Error " + err.Error())
		os.Exit(1)
	}

	feedRepository := repository.NewFeedRepository(dbConnection)
	usersRepository := repository.NewUsersRepository(dbConnection)
	loggerService := &logger.Logger{}
	eventManager := events.NewEventManager(loggerService)
	feedParser := parser.NewParser(&http.Client{})
	registerHandlers(eventManager, telegramClient, loggerService, usersRepository, feedRepository, feedParser)
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
		//feed.NextParse = nowTimestamp + ParsePeriod
		//for _, feed := range feeds {
		//	go parser.ParseFeed(feed)
		//}

		time.Sleep(2 * time.Second)
	}
}

func registerHandlers(
	em *events.Manager,
	messenger *telegram.Client,
	logger *logger.Logger,
	usersRepository *repository.UsersRepository,
	feedRepository *repository.FeedRepository,
	feedParser *parser.Parser,
) {
	newFeedItemHandler := events.NewNewFeedItemHandler(messenger, logger, usersRepository)
	em.RegisterHandler(newFeedItemHandler)

	newMessageHandler := events.NewNewMessageHandler(logger, usersRepository, messenger, em)
	em.RegisterHandler(newMessageHandler)

	startChatHandler := events.NewStartChatHandler(messenger, logger, usersRepository)
	em.RegisterHandler(startChatHandler)

	receiveAddMessageHandler := events.NewReceiveAddMessageHandler(messenger, logger, usersRepository)
	em.RegisterHandler(receiveAddMessageHandler)

	receiveDeleteMessageHandler := events.NewReceiveDeleteMessageHandler(
		messenger,
		logger,
		usersRepository,
		feedRepository,
	)
	em.RegisterHandler(receiveDeleteMessageHandler)

	addFeedHandler := events.NewAddFeedHandler(
		messenger,
		logger,
		usersRepository,
		feedRepository,
		feedParser,
		em,
	)
	em.RegisterHandler(addFeedHandler)

	deleteFeedHandler := events.NewDeleteFeedHandler(
		messenger,
		logger,
		usersRepository,
		feedRepository,
		feedParser,
		em,
	)
	em.RegisterHandler(deleteFeedHandler)
}
