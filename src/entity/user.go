package entity

type User struct {
	Id          int    `db:"id"`
	TelegramId  int    `db:"telegram_id"`
	ChatId      int    `db:"chat_id"`
	LastMessage string `db:"last_message"`
}

func NewUser(telegramId int, chatId int) *User {
	return &User{
		TelegramId: telegramId,
		ChatId:     chatId,
	}
}
