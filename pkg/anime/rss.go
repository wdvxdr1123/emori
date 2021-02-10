package anime

import (
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
)

type annieTorrent struct {
	title       string
	description string
	date        *time.Time
	link        string
	torrentLink string
}

func getAnnieInfo(t []Tag) ([]annieTorrent, error) {
	var tagID []string
	for i := range t {
		tagID = append(tagID, t[i].ID)
	}
	rssURL := "https://bangumi.moe/rss/tags/" + strings.Join(tagID, "+")
	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(rssURL)
	if err != nil {
		return nil, errors.Wrap(err, "parse rss failed")
	}
	var ret []annieTorrent
	for _, item := range feed.Items {
		ret = append(ret, annieTorrent{
			title:       item.Title,
			description: item.Description,
			date:        item.PublishedParsed,
			link:        item.Link,
			torrentLink: item.Enclosures[0].URL,
		})
	}
	return ret, nil
}
