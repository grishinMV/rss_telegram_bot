package events

import (
	"errors"
	"reflect"
	"rss-bot/src/entity"
	"rss-bot/src/repository"
	"rss-bot/src/telegram"
)

type ReceiveAddMessage struct {
	User    *entity.User
	Message *telegram.Message
}

func (fu ReceiveAddMessage) GetName() string {
	return "ReceiveAddMessage"
}

type ReceiveAddMessageHandler struct {
	messenger       *telegram.Client
	logger          Logger
	usersRepository *repository.UsersRepository
}

func NewReceiveAddMessageHandler(
	messenger *telegram.Client,
	logger Logger,
	usersRepository *repository.UsersRepository,
) *ReceiveAddMessageHandler {
	return &ReceiveAddMessageHandler{
		messenger:       messenger,
		logger:          logger,
		usersRepository: usersRepository,
	}
}

func (h *ReceiveAddMessageHandler) GetEventName() string {
	return "ReceiveAddMessage"
}

func (h *ReceiveAddMessageHandler) Handle(e interface{}) error {
	event, ok := e.(ReceiveAddMessage)
	if ok == false {
		t := reflect.TypeOf(e)
		return errors.New("Incorrect type assertion ReceiveAddMessage " + t.String())
	}

	user := event.User
	message := event.Message

	user.LastMessage = message.Text
	err := h.usersRepository.Save(user)
	_, err = h.messenger.SendTextMessage(user.ChatId, "Отправь ссылку на rss")
	if err != nil {
		h.logger.Log(err.Error())
		return err
	}

	return nil
}
