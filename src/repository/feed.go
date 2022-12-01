package repository

import (
	"rss-bot/src/db"
	"rss-bot/src/entity"
	"strconv"
)

type FeedRepository struct {
	conn *db.Connection
}

func NewFeedRepository(conn *db.Connection) *FeedRepository {
	return &FeedRepository{conn: conn}
}

func (r *FeedRepository) FindByUser(userId int) ([]entity.Feed, error) {
	query := "select f.* from feeds f inner join users_feeds uf on f.id = uf.feed_id where uf.user_id = " + strconv.Itoa(userId)
	var feed []entity.Feed
	err := r.conn.Db.Select(&feed, query)
	if err != nil {
		return nil, err
	}

	return feed, nil
}

func (r *FeedRepository) FindForUpdate() ([]entity.Feed, error) {
	query := "select * from rss_parser.feeds where feeds.next_parse < CURRENT_TIMESTAMP()"
	var feed []entity.Feed
	err := r.conn.Db.Select(&feed, query)
	if err != nil {
		return nil, err
	}

	return feed, nil
}

func (r *FeedRepository) Update(f entity.Feed) error {
	query := "UPDATE feeds SET link = ?, last_new = ?, next_parse = ? WHERE id = ?"
	_, err := r.conn.Db.Exec(query, f.Link, f.LastNew, f.NextParse, f.Id)
	if err != nil {
		return err
	}

	return nil
}
