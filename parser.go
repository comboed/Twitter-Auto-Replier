package main

import (
	"github.com/valyala/fastjson"
	"strconv"
	"strings"
)

func parseTweet(body string) (string, string) {
	var json, _ = fastjson.Parse(body)
	var timeline *fastjson.Value = json.Get("data", "user", "result", "timeline_v2", "timeline", "instructions")
	var entries *fastjson.Value = getEntries(timeline)
	var tweetID, tweetDate string = getTweet(entries)
	
	return tweetID, tweetDate
}

func getEntries(timeline *fastjson.Value) *fastjson.Value {
	for i := 0; i < 20 && timeline != nil; i++ {
		var entries *fastjson.Value = timeline.Get(strconv.Itoa(i), "entries")
		if entries != nil {
			return entries
		}
	}
	return nil
}

func getTweet(entries *fastjson.Value) (string, string) {
	for i := 0; i < 20 && entries != nil; i++ {
		var index string = strconv.Itoa(i)
		if !strings.Contains(entries.Get(index, "entryId").String(), "promoted-tweet") {
			var singleTweet *fastjson.Value = entries.Get(index, "content", "itemContent", "tweet_results", "result")
			var multiTweet *fastjson.Value = entries.Get(index, "content", "items", "0", "item", "itemContent", "tweet_results", "result")

			if singleTweet != nil && singleTweet.Get("legacy", "retweeted_status_result") == nil {
				return string(singleTweet.GetStringBytes("rest_id")), string(singleTweet.GetStringBytes("legacy", "created_at"))
			} else if multiTweet != nil && multiTweet.Get("legacy", "retweeted_status_result") == nil {
				return string(multiTweet.GetStringBytes("rest_id")), string(multiTweet.GetStringBytes("legacy", "created_at"))
			}
		}
	}
	return "", ""
}

