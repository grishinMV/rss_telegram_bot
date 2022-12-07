package parser

import (
	"github.com/mmcdole/gofeed"
	"net/http"
	"rss-bot/src/entity"
	"time"
)

const userAgent = "RSS bot"

var ParsePeriod = int64(60 * time.Second)

type Parser struct {
	httpClient *http.Client
	parser     *gofeed.Parser
}

func NewParser(
	client *http.Client,
) *Parser {
	return &Parser{
		httpClient: client,
		parser:     gofeed.NewParser(),
	}
}

func (p *Parser) Parse(feed *entity.Feed) ([]FeedItem, error) {
	nowTimestamp := time.Now().Unix()

	req, err := http.NewRequest("GET", feed.Link, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)

	response, err := p.httpClient.Do(req)
	if err != nil {
		feed.LastNew = nowTimestamp

		return nil, err
	}

	parsedRss, err := p.parser.Parse(response.Body)
	if err != nil {
		return nil, err
	}

	var newItems []FeedItem
	for _, item := range parsedRss.Items {
		if item.PublishedParsed.Unix() < feed.LastNew {
			continue
		}

		newItems = append(newItems, FeedItem{
			FeedId: feed.Id,
			Text:   item.Title,
			Link:   item.Link,
		})
	}

	feed.LastNew = nowTimestamp
	feed.NextParse = ParsePeriod + time.Now().Unix()

	return newItems, err
}
