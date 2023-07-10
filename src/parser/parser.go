package parser

import (
	"github.com/mmcdole/gofeed"
	"log"
	"net/http"
	"rss-bot/src/entity"
	"rss-bot/src/logger"
	"time"
)

const userAgent = "RSS bot"

var ParsePeriod int64 = 60

type Parser struct {
	httpClient http.Client
	parser     gofeed.Parser
	logger     *logger.Logger
}

func NewParser(
	client http.Client,
	logger *logger.Logger,
) *Parser {
	return &Parser{
		httpClient: client,
		parser:     *gofeed.NewParser(),
		logger:     logger,
	}
}

func (p *Parser) Parse(feed *entity.Feed, location *time.Location) ([]FeedItem, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
		}
	}()

	nowTimestamp := time.Now().In(location).Unix()

	req, err := http.NewRequest("GET", feed.Link, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)

	response, err := p.httpClient.Do(req)
	if err != nil {
		feed.NextParse = ParsePeriod + time.Now().In(location).Unix()

		return nil, err
	}

	parsedRss, err := p.parser.Parse(response.Body)
	if err != nil {
		return nil, err
	}

	var newItems []FeedItem

	for _, item := range parsedRss.Items {
		if item == nil {
			continue
		}

		if item.PublishedParsed == nil {
			p.logger.Log(item.Link + " [skip] PublishedParsed is null")

			continue
		}

		if item.PublishedParsed.In(location).Unix() < feed.LastNew {
			p.logger.Log(item.Link + " [skip] item time < LastNew")

			continue
		}

		newItems = append(newItems, FeedItem{
			FeedId: feed.Id,
			Text:   item.Title,
			Link:   item.Link,
		})
	}

	feed.LastNew = nowTimestamp
	feed.NextParse = ParsePeriod + nowTimestamp

	return newItems, err
}
