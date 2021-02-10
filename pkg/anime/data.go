package anime

import (
	"strconv"

	jsoniter "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

var json = jsoniter.ConfigFastest

var db *leveldb.DB

func init() {
	var err error
	db, err = leveldb.OpenFile("anime.db", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func queryTags(userID int64) []Tag {
	var tags = make([]Tag, 0)
	val, _ := db.Get(helper.StringToBytes(strconv.FormatInt(userID, 10)), nil)
	_ = json.Unmarshal(val, &tags)
	return tags
}

func addTags(userID int64, tag Tag) {
	var tags = queryTags(userID)
	for i := range tags {
		if tags[i].ID == tag.ID {
			return
		}
	}
	tags = append(tags, tag)
	newVal, _ := json.Marshal(&tags)
	_ = db.Put([]byte(strconv.FormatInt(userID, 10)), newVal, nil)
}

func deleteTags(userID int64, tag Tag) {
	var tags = queryTags(userID)
	tot := len(tags)
	for i := range tags {
		if tags[i].ID == tag.ID {
			tags = append(tags[:i], tags[i+1:tot]...)
			newVal, _ := json.Marshal(&tags)
			_ = db.Put([]byte(strconv.FormatInt(userID, 10)), newVal, nil)
			return
		}
	}
}
