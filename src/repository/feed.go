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

func (r *FeedRepository) Delete(feed *entity.Feed) error {
	query := "DELETE FROM feeds WHERE id = ?"
	_, err := r.conn.Db.Exec(query, feed.Id)
	if err != nil {
		return err
	}

	return nil
}

func (r *FeedRepository) GetRelationsCount(feed *entity.Feed) (int, error) {
	query := "SELECT count(*) AS count FROM users_feeds WHERE feed_id = ?"

	var count int

	err := r.conn.Db.Get(&count, query, feed.Id)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *FeedRepository) Save(feed *entity.Feed) error {
	var err error
	if feed.Id == 0 {
		query := "INSERT INTO feeds (link, last_new, next_parse) VALUES (?, ?, ?);"
		result, err := r.conn.Db.Exec(query, feed.Link, feed.LastNew, feed.NextParse)

		feedId, err := result.LastInsertId()
		if err != nil {
			return err
		}

		feed.Id = int(feedId)
	} else {
		query := "UPDATE feeds SET link = ?, last_new = ?, next_parse = ? WHERE id = ?;"

		_, err = r.conn.Db.Exec(query, feed.Link, feed.LastNew, feed.NextParse, feed.Id)
	}

	return err
}

func (r *FeedRepository) FindByUser(userId int) ([]entity.Feed, error) {
	query := "SELECT f.* FROM feeds f INNER JOIN users_feeds uf ON f.id = uf.feed_id WHERE uf.user_id = ?"
	var feed []entity.Feed
	err := r.conn.Db.Select(&feed, query, strconv.Itoa(userId))
	if err != nil {
		return nil, err
	}

	return feed, nil
}

func (r *FeedRepository) FindByUrl(url string) ([]entity.Feed, error) {
	query := "SELECT * FROM feeds WHERE link = ?"
	var feeds []entity.Feed
	err := r.conn.Db.Select(&feeds, query, url)
	if err != nil {
		return nil, err
	}

	return feeds, nil
}

func (r *FeedRepository) FindForUpdate() ([]entity.Feed, error) {
	query := "SELECT * FROM feeds WHERE feeds.next_parse < CURRENT_TIMESTAMP()"
	var feed []entity.Feed
	err := r.conn.Db.Select(&feed, query)
	if err != nil {
		return nil, err
	}

	return feed, nil
}

func (r *FeedRepository) Update(f *entity.Feed) error {
	query := "UPDATE feeds SET link = ?, last_new = ?, next_parse = ? WHERE id = ?"
	_, err := r.conn.Db.Exec(query, f.Link, f.LastNew, f.NextParse, f.Id)
	if err != nil {
		return err
	}

	return nil
}
