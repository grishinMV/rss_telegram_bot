package events

import (
	"errors"
	"reflect"
	"rss-bot/src/entity"
	"rss-bot/src/repository"
	"rss-bot/src/telegram"
)

type ReceiveDeleteMessage struct {
	User    *entity.User
	Message *telegram.Message
}

func (fu ReceiveDeleteMessage) GetName() string {
	return "ReceiveDeleteMessage"
}

type ReceiveDeleteMessageHandler struct {
	messenger       *telegram.Client
	logger          Logger
	usersRepository *repository.UsersRepository
	feedsRepository *repository.FeedRepository
}

func NewReceiveDeleteMessageHandler(
	messenger *telegram.Client,
	logger Logger,
	usersRepository *repository.UsersRepository,
	feedsRepository *repository.FeedRepository,
) *ReceiveDeleteMessageHandler {
	return &ReceiveDeleteMessageHandler{
		messenger:       messenger,
		logger:          logger,
		usersRepository: usersRepository,
		feedsRepository: feedsRepository,
	}
}

func (h *ReceiveDeleteMessageHandler) GetEventName() string {
	return "ReceiveDeleteMessage"
}

func (h *ReceiveDeleteMessageHandler) Handle(e interface{}) error {
	event, ok := e.(ReceiveDeleteMessage)
	if ok == false {
		t := reflect.TypeOf(e)
		return errors.New("Incorrect type assertion ReceiveDeleteMessage " + t.String())
	}

	user := event.User
	message := event.Message

	user.LastMessage = message.Text
	err := h.usersRepository.Save(user)

	feeds, err := h.feedsRepository.FindByUser(event.User.Id)
	if err != nil {
		return err
	}
	var buttons [][]telegram.KeyboardButton

	for _, feed := range feeds {
		buttons = append(buttons, []telegram.KeyboardButton{{Text: feed.Link}})
	}

	_, err = h.messenger.SendReplyMarkup(
		event.Message.Chat.Id,
		"Выберите что удалить",
		telegram.ReplyKeyboardMarkup{
			Keyboard:        buttons,
			OneTimeKeyboard: true,
		})

	if err != nil {
		h.logger.Log(err.Error())
		return err
	}

	return nil
}
