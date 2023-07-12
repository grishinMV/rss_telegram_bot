package events

import (
	"database/sql"
	"errors"
	"reflect"
	"regexp"
	"rss-bot/src/entity"
	"rss-bot/src/repository"
	"rss-bot/src/telegram"
	"strconv"
)

const UserMessageStart = "/start"
const UserMessageAdd = "/add"
const UserMessageDelete = "/delete"

type commands struct {
	Start, Add, Delete *regexp.Regexp
}

type NewMessage struct {
	Message *telegram.Message
}

func (m NewMessage) GetName() string {
	return "NewMessage"
}

type NewMessageHandler struct {
	logger          Logger
	usersRepository *repository.UsersRepository
	telegram        telegram.Client
	em              *Manager
	commands        commands
}

func NewNewMessageHandler(
	logger Logger,
	usersRepository *repository.UsersRepository,
	telegram telegram.Client,
	em *Manager,
) *NewMessageHandler {
	return &NewMessageHandler{
		logger:          logger,
		usersRepository: usersRepository,
		telegram:        telegram,
		em:              em,
		commands: commands{
			Start:  regexp.MustCompile(UserMessageStart),
			Add:    regexp.MustCompile(UserMessageAdd),
			Delete: regexp.MustCompile(UserMessageDelete),
		},
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

	switch {
	case h.commands.Start.MatchString(event.Message.Text):
		go h.em.Dispatch(StartChat{
			Message: event.Message,
			User:    user,
		})

		return nil
	case h.commands.Add.MatchString(event.Message.Text):
		go h.em.Dispatch(ReceiveAddMessage{
			Message: event.Message,
			User:    user,
		})

		return nil
	case h.commands.Delete.MatchString(event.Message.Text):
		go h.em.Dispatch(ReceiveDeleteMessage{
			Message: event.Message,
			User:    user,
		})

		return nil
	default:
		return h.handleCustomMessage(user, event.Message)
	}
}

func (h *NewMessageHandler) handleCustomMessage(user *entity.User, message *telegram.Message) error {
	switch {
	case h.commands.Add.MatchString(user.LastMessage):
		go h.em.Dispatch(AddFeed{
			User:    user,
			Message: message,
		})
	case h.commands.Delete.MatchString(user.LastMessage):
		go h.em.Dispatch(DeleteFeed{
			User:    user,
			Message: message,
		})
	default:
		_, _ = h.telegram.SendTextMessage(message.Chat.Id, "Что?")
	}

	user.LastMessage = message.Text
	return h.usersRepository.Save(user)
}
