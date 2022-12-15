package events

import (
	"errors"
	"reflect"
	"rss-bot/src/entity"
	"rss-bot/src/repository"
	"rss-bot/src/telegram"
)

type StartChat struct {
	User    *entity.User
	Message *telegram.Message
}

func (fu StartChat) GetName() string {
	return "StartChat"
}

type StartChatHandler struct {
	messenger       telegram.Client
	logger          Logger
	usersRepository *repository.UsersRepository
}

func NewStartChatHandler(
	messenger telegram.Client,
	logger Logger,
	usersRepository *repository.UsersRepository,
) *StartChatHandler {
	return &StartChatHandler{
		messenger:       messenger,
		logger:          logger,
		usersRepository: usersRepository,
	}
}

func (h *StartChatHandler) GetEventName() string {
	return "StartChat"
}

func (h *StartChatHandler) Handle(e interface{}) error {
	event, ok := e.(StartChat)
	if ok == false {
		t := reflect.TypeOf(e)
		return errors.New("Incorrect type assertion StartChat " + t.String())
	}

	user := event.User
	message := event.Message

	if event.User == nil {
		user = entity.NewUser(message.From.Id, message.Chat.Id)
		user.LastMessage = message.Text
		err := h.usersRepository.Save(user)
		if err != nil {
			h.logger.Log(err.Error())
			return err
		}
	}

	_, err := h.messenger.SendTextMessage(user.ChatId, "Привет! Ты можешь добавить новую рассылку командой /add")

	return err
}
