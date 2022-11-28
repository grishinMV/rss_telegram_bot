package repository

import (
	"rss-bot/src/db"
	"rss-bot/src/entity"
)

type UsersRepository struct {
	conn *db.Connection
}

func NewUsersRepository(conn *db.Connection) *UsersRepository {
	return &UsersRepository{conn: conn}
}

func (r *UsersRepository) FindUsersByFeedId(feedId int) ([]entity.User, error) {
	query := "select u.* from users u inner join users_feeds uf on u.id = uf.user_id where uf.feed_id = ?"

	var users []entity.User
	err := r.conn.Db.Select(&users, query, feedId)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UsersRepository) FindByTelegramId(telegramId int) (*entity.User, error) {
	query := "select u.* from users u where u.telegram_id = ?"

	var user entity.User
	err := r.conn.Db.Get(&user, query, telegramId)
	if err != nil {
		return nil, err
	}

	return &user, err
}

func (r *UsersRepository) Save(user *entity.User) error {
	var err error
	if user.Id == 0 {
		query := "INSERT INTO users (telegram_id, chat_id, last_message) VALUES (?, ?, ?);"
		_, err = r.conn.Db.Exec(query, user.ChatId, user.ChatId, user.LastMessage)
	} else {
		query := "UPDATE users SET telegram_id = ?, chat_id = ?, last_message = ? WHERE id = ?;"

		_, err = r.conn.Db.Exec(query, user.TelegramId, user.ChatId, user.LastMessage, user.Id)
	}

	return err
}
