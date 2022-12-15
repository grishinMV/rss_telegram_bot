package events

import (
	"errors"
	"net/url"
	"reflect"
	"rss-bot/src/entity"
	"rss-bot/src/parser"
	"rss-bot/src/repository"
	"rss-bot/src/telegram"
	"strings"
	"time"
)

type AddFeed struct {
	User    *entity.User
	Message *telegram.Message
}

func (fu AddFeed) GetName() string {
	return "AddFeed"
}

type AddFeedHandler struct {
	messenger       telegram.Client
	logger          Logger
	usersRepository *repository.UsersRepository
	feedsRepository *repository.FeedRepository
	feedParser      *parser.Parser
	em              *Manager
	location        *time.Location
}

func NewAddFeedHandler(
	messenger telegram.Client,
	logger Logger,
	usersRepository *repository.UsersRepository,
	feedsRepository *repository.FeedRepository,
	feedParser *parser.Parser,
	em *Manager,
	location *time.Location,
) *AddFeedHandler {
	return &AddFeedHandler{
		messenger:       messenger,
		logger:          logger,
		usersRepository: usersRepository,
		feedsRepository: feedsRepository,
		feedParser:      feedParser,
		em:              em,
		location:        location,
	}
}

func (h *AddFeedHandler) GetEventName() string {
	return "AddFeed"
}

func (h *AddFeedHandler) Handle(e interface{}) error {
	event, ok := e.(AddFeed)
	if ok == false {
		t := reflect.TypeOf(e)
		return errors.New("Incorrect type assertion AddFeed " + t.String())
	}

	link, err := h.prepareLink(event.Message.Text)
	if err != nil {
		h.handleError(event.Message)
		return nil
	}

	feeds, err := h.feedsRepository.FindByUrl(link)
	var feed entity.Feed

	if len(feeds) > 0 {
		feed = feeds[0]
	} else {
		feed = entity.Feed{
			Link:    link,
			LastNew: time.Now().In(h.location).Unix(),
		}

		_, err = h.feedParser.Parse(&feed, h.location)
		if err != nil {
			h.handleError(event.Message)
			return nil
		}

		err = h.feedsRepository.Save(&feed)
		if err != nil {
			h.logger.Log(err.Error())
			_, err = h.messenger.SendTextMessage(event.Message.Chat.Id, "Ошибка")
			return err
		}
	}

	err = h.usersRepository.AddFeed(event.User, feed)
	if err != nil {
		h.logger.Log(err.Error())
		_, err = h.messenger.SendTextMessage(event.Message.Chat.Id, "Ошибка. Возможно вы уже добавили такую ссылку")
		return err
	}

	_, err = h.messenger.SendTextMessage(event.Message.Chat.Id, "Успешно")
	if err != nil {
		h.logger.Log(err.Error())
		return err
	}

	return nil
}

func (h *AddFeedHandler) prepareLink(feedUrl string) (string, error) {
	feedUrl = strings.Trim(feedUrl, " ")
	feedUrl = strings.Trim(feedUrl, "\n")
	parsedUrl, err := url.Parse(feedUrl)
	if err != nil {
		return "", nil
	}

	return parsedUrl.String(), nil
}

func (h *AddFeedHandler) handleError(message *telegram.Message) {
	_, _ = h.messenger.SendTextMessage(
		message.Chat.Id,
		"Не получилось добавить url. Воможно кривая ссылка",
	)
}
