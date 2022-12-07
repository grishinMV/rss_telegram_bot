package events

import (
	"errors"
	"reflect"
	"rss-bot/src/entity"
	"rss-bot/src/parser"
	"rss-bot/src/repository"
	"rss-bot/src/telegram"
)

type DeleteFeed struct {
	User    *entity.User
	Message *telegram.Message
}

func (fu DeleteFeed) GetName() string {
	return "DeleteFeed"
}

type DeleteFeedHandler struct {
	messenger       *telegram.Client
	logger          Logger
	usersRepository *repository.UsersRepository
	feedsRepository *repository.FeedRepository
	feedParser      *parser.Parser
	em              *Manager
}

func NewDeleteFeedHandler(
	messenger *telegram.Client,
	logger Logger,
	usersRepository *repository.UsersRepository,
	feedsRepository *repository.FeedRepository,
	feedParser *parser.Parser,
	em *Manager,
) *DeleteFeedHandler {
	return &DeleteFeedHandler{
		messenger:       messenger,
		logger:          logger,
		usersRepository: usersRepository,
		feedsRepository: feedsRepository,
		feedParser:      feedParser,
		em:              em,
	}
}

func (h *DeleteFeedHandler) GetEventName() string {
	return "DeleteFeed"
}

func (h *DeleteFeedHandler) Handle(e interface{}) error {
	event, ok := e.(DeleteFeed)
	if ok == false {
		t := reflect.TypeOf(e)
		return errors.New("Incorrect type assertion DeleteFeed " + t.String())
	}

	feeds, err := h.feedsRepository.FindByUrl(event.Message.Text)
	if err != nil {
		_, _ = h.messenger.SendTextMessage(event.Message.Chat.Id, "Произошла ошибка при удалении")
		return err
	}
	for _, feed := range feeds {
		err := h.usersRepository.DeleteFeed(event.User, feed)
		if err != nil {
			_, _ = h.messenger.SendTextMessage(event.Message.Chat.Id, "Произошла ошибка при удалении")
			return err
		}

		count, err := h.feedsRepository.GetRelationsCount(&feed)
		if err != nil {
			return err
		}

		if count == 0 {
			err := h.feedsRepository.Delete(&feed)
			if err != nil {
				return err
			}
		}
	}
	_, _ = h.messenger.SendTextMessage(event.Message.Chat.Id, "Успешно")

	return nil
}
