package controller

import (
	"fmt"
	"strings"

	"github.com/nexryai/polyxia/server/feed"
)

const realtimeFeedUrl = "https://www.data.jma.go.jp/developer/xml/feed/eqvol.xml"
const longtermFeedUrl = "https://www.data.jma.go.jp/developer/xml/feed/eqvol_l.xml"

func GetJMAEvents() (*[]string, error) {
	realtimeFeed, err := feed.FetchFeed(realtimeFeedUrl)
	if err != nil {
		fmt.Printf("Failed to fetch feed: %v\n", err)
		return nil, err
	}

	longtermFeed, err := feed.FetchFeed(longtermFeedUrl)
	if err != nil {
		fmt.Printf("Failed to fetch feed: %v\n", err)
		return nil, err
	}

	// イベントの詳細情報のXMLへのリンクを格納する配列
	var events []string

	// リアルタイムフィードから追加
	for _, item := range realtimeFeed.Items {
		// 対応してる情報のみを扱う
		if strings.Contains(item.Link, "_VXSE51_") || strings.Contains(item.Link, "_VXSE53_") || strings.Contains(item.Link, "_VTSE41") {
			// ID抽出
			e := strings.TrimSuffix(strings.TrimPrefix(item.Link, "https://www.data.jma.go.jp/developer/xml/data/"), ".xml")
			events = append(events, e)
		}
	}

	// 長期フィードから追加
	for _, item := range longtermFeed.Items {
		if strings.Contains(item.Link, "_VXSE51_") || strings.Contains(item.Link, "_VXSE53_") || strings.Contains(item.Link, "_VTSE41") {
			// ID抽出
			e := strings.TrimSuffix(strings.TrimPrefix(item.Link, "https://www.data.jma.go.jp/developer/xml/data/"), ".xml")

			// 既に存在するイベントならパス
			for _, event := range events {
				if event == e {
					continue
				}
			}

			events = append(events, e)
		}
	}

	return &events, nil
}
