package events

import (
	"database/sql"
	"errors"
	"reflect"
	"rss-bot/src/entity"
	"rss-bot/src/repository"
	"rss-bot/src/telegram"
	"strconv"
)

const UserMessageStart = "/start"
const UserMessageAdd = "/add"
const UserMessageDelete = "/delete"

type NewMessage struct {
	Message *telegram.Message
}

func (m NewMessage) GetName() string {
	return "NewMessage"
}

type NewMessageHandler struct {
	logger          Logger
	usersRepository *repository.UsersRepository
	telegram        *telegram.Client
}

func NewNewMessageHandler(
	logger Logger,
	usersRepository *repository.UsersRepository,
	telegram *telegram.Client,
) *NewMessageHandler {
	return &NewMessageHandler{
		logger:          logger,
		usersRepository: usersRepository,
		telegram:        telegram,
	}
}

func (h *NewMessageHandler) GetEventName() string {
	return "NewMessage"
}

func (h *NewMessageHandler) Handle(e interface{}) error {
	event, ok := e.(NewMessage)
	if ok == false {
		t := reflect.TypeOf(e)
		return errors.New("Incorrect type assertion NewMessage " + t.String())
	}

	user, err := h.usersRepository.FindByTelegramId(event.Message.From.Id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			h.logger.Log("telegaramUserId = " + strconv.Itoa(event.Message.From.Id) + " _ " + err.Error())
			return err
		}
	}

	if user == nil {
		return h.handleStartMessage(user, event.Message)
	}

	switch event.Message.Text {
	case UserMessageStart:
		return h.handleStartMessage(user, event.Message)
	case UserMessageAdd:
		return h.handleAddMessage(user, event.Message)
	case UserMessageDelete:
		return h.handleDeleteMessage(user, event.Message)
	default:
		return h.handleCustomMessage(user, event.Message)
	}
}

func (h *NewMessageHandler) handleStartMessage(user *entity.User, message *telegram.Message) error {
	if user == nil {
		user = entity.NewUser(message.From.Id, message.Chat.Id)
		user.LastMessage = message.Text
		err := h.usersRepository.Save(user)
		if err != nil {
			h.logger.Log(err.Error())
			return err
		}
	}
	_, _ = h.telegram.SendTextMessage(user.ChatId, "Привет! Ты можешь добавить новую рассылку командой /add")

	return nil
}

func (h *NewMessageHandler) handleAddMessage(user *entity.User, message *telegram.Message) error {
	user.LastMessage = message.Text
	err := h.usersRepository.Save(user)
	_, err = h.telegram.SendTextMessage(user.ChatId, "Отправь ссылку на rss")
	if err != nil {
		h.logger.Log(err.Error())
		return err
	}

	return nil
}

func (h *NewMessageHandler) handleDeleteMessage(user *entity.User, message *telegram.Message) error {
	user.LastMessage = message.Text
	err := h.usersRepository.Save(user)
	_, err = h.telegram.SendTextMessage(user.ChatId, "Я пока не умею удалять подписки, но обязательно научусь🗿")
	if err != nil {
		h.logger.Log(err.Error())
		return err
	}

	return nil
}

func (h *NewMessageHandler) handleCustomMessage(user *entity.User, message *telegram.Message) error {
	switch user.LastMessage {
	case UserMessageAdd:
		_, _ = h.telegram.SendTextMessage(message.Chat.Id, "add")
	case UserMessageDelete:
		h.logger.Log("delete")
		_, _ = h.telegram.SendTextMessage(message.Chat.Id, "delete")
	}

	user.LastMessage = message.Text
	return h.usersRepository.Save(user)
}
