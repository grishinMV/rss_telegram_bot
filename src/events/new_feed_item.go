package events

import (
	"errors"
	"net/url"
	"reflect"
	"rss-bot/src/repository"
	"rss-bot/src/telegram"
)

type NewFeedItem struct {
	FeedId int
	Text   string
	Link   string
}

func (fu NewFeedItem) GetName() string {
	return "NewFeedItem"
}

type NewFeedItemHandler struct {
	messenger       *telegram.Client
	logger          Logger
	usersRepository *repository.UsersRepository
}

func NewNewFeedItemHandler(
	messenger *telegram.Client,
	logger Logger,
	usersRepository *repository.UsersRepository,
) *NewFeedItemHandler {
	return &NewFeedItemHandler{
		messenger:       messenger,
		logger:          logger,
		usersRepository: usersRepository,
	}

}

func (h *NewFeedItemHandler) GetEventName() string {
	return "NewFeedItem"
}

func (h *NewFeedItemHandler) Handle(e interface{}) error {
	event, ok := e.(NewFeedItem)
	if ok == false {
		t := reflect.TypeOf(e)
		return errors.New("Incorrect type assertion NewFeedItem " + t.String())
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
