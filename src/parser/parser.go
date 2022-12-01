package parser

import (
	"github.com/mmcdole/gofeed"
	"net/http"
	"rss-bot/src/entity"
	"rss-bot/src/events"
	"rss-bot/src/repository"
	"time"
)

const ParsePeriod = 600 // 10 минут
const userAgent = "RSS notify bot t.me/rss_feed_gmv_bot"

type Parser struct {
	httpClient     *http.Client
	parser         *gofeed.Parser
	feedRepository *repository.FeedRepository
	EventManager   *events.Manager
}

func NewParser(
	eventManager *events.Manager,
	feedRepository *repository.FeedRepository,
	client *http.Client,
) Parser {
	return Parser{
		EventManager:   eventManager,
		httpClient:     client,
		parser:         gofeed.NewParser(),
		feedRepository: feedRepository,
	}
}

func (p *Parser) ParseFeed(feed entity.Feed) {
	nowTimestamp := time.Now().Unix()
	err := p.feedRepository.Update(feed)
	if err != nil {
		return
	}

	feed.NextParse = nowTimestamp + ParsePeriod
	if feed.LastNew+ParsePeriod > nowTimestamp {
		return
	}

	req, err := http.NewRequest("GET", feed.Link, nil)
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", userAgent)

	response, err := p.httpClient.Do(req)
	if err != nil {
		feed.LastNew = nowTimestamp
		err = p.feedRepository.Update(feed)
		if err != nil {
			return
		}

		return
	}

	parsedRss, err := p.parser.Parse(response.Body)
	if err != nil {
		return
	}

	for _, item := range parsedRss.Items {
		if item.PublishedParsed.Unix() < feed.LastNew {
			continue
		}

		go p.EventManager.Dispatch(events.FeedUpdated{
			FeedId: feed.Id,
			Text:   item.Title,
			Link:   item.Link,
		})
	}

	feed.LastNew = nowTimestamp
	err = p.feedRepository.Update(feed)
	if err != nil {
		return
	}

	return
}
