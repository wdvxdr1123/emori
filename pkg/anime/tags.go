package anime

import (
	"github.com/levigross/grequests"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

type Tag struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

func GetTagByKeyword(keyword string) ([]Tag, error) {
	ro := &grequests.RequestOptions{
		JSON: map[string]interface{}{
			"keywords": true,
			"multi":    true,
			"name":     keyword,
		},
	}
	resp, err := grequests.Post("https://bangumi.moe/api/tag/search", ro)
	if err != nil {
		return nil, errors.Wrap(err, "get tag failed")
	}
	data := gjson.ParseBytes(resp.Bytes())
	var ret []Tag
	data.Get("tag").ForEach(func(_, v gjson.Result) bool {
		ret = append(ret, Tag{
			ID:   v.Get("_id").Str,
			Type: v.Get("type").Str,
			Name: func() string {
				locale := v.Get("locale")
				switch {
				case locale.Get("zh_cn").Str != "":
					return locale.Get("zh_cn").Str
				case locale.Get("zh_tw").Str != "":
					return locale.Get("zh_tw").Str
				case locale.Get("ja").Str != "":
					return locale.Get("ja").Str
				case locale.Get("en").Str != "":
					return locale.Get("en").Str
				default:
					return "<unknown>"
				}
			}(),
		})
		return true
	})
	return ret, nil
}
