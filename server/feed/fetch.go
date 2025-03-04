package feed

import (
	"errors"
	"net/http"

	"github.com/mmcdole/gofeed"
	"github.com/nexryai/archer"
)

func FetchFeed(url string) (*gofeed.Feed, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	requester := archer.SecureRequest{
		Request:     req,
		TimeoutSecs: 10,
		MaxSize:     1024 * 1024 * 10,
	}

	resp, err := requester.Send()
	if err != nil {
		return nil, err
	}

	parser := gofeed.NewParser()
	feed, err := parser.Parse(resp.Body)
	if err != nil {
		return nil, err
	} else if feed == nil {
		return nil, errors.New("feed is nil")
	}

	return feed, nil
}