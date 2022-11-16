package events

import (
	"errors"
	"net/url"
	"reflect"
	"rss-bot/src/repository"
	"rss-bot/src/telegram"
)

type FeedUpdated struct {
	FeedId int
	Text   string
	Link   string
}

func (fu FeedUpdated) GetName() string {
	return "FeedUpdated"
}

type FeedUpdatedHandler struct {
	messenger       *telegram.Client
	logger          Logger
	usersRepository *repository.UsersRepository
}

func NewFeedUpdatedHandler(
	messenger *telegram.Client,
	logger Logger,
	usersRepository *repository.UsersRepository,
) *FeedUpdatedHandler {
	return &FeedUpdatedHandler{
		messenger:       messenger,
		logger:          logger,
		usersRepository: usersRepository,
	}

}

func (h *FeedUpdatedHandler) GetEventName() string {
	return "FeedUpdated"
}

func (h *FeedUpdatedHandler) Handle(e interface{}) error {
	event, ok := e.(FeedUpdated)
	if ok == false {
		t := reflect.TypeOf(e)
		return errors.New("Incorrect type assertion FeedUpdated " + t.String())
	}

	users, err := h.usersRepository.FindUsersByFeedId(event.FeedId)
	if err != nil {
		return err
	}

	for _, user := range users {
		_, err := h.messenger.SendTextMessage(user.ChatId, event.Text+" "+url.QueryEscape(event.Link))
		if err != nil {
			h.logger.Log(err.Error())
		}
	}

	return nil
}
