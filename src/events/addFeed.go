package events

import (
	"rss-bot/src/entity"
)

type AddFeed struct {
	user *entity.User
	link string
}

func (m AddFeed) GetName() string {
	return "AddFeed"
}
