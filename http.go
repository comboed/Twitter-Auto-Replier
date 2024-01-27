package main

import (
	"github.com/valyala/fasthttp"
	"github.com/valyala/fastjson"
	"strings"
)


func isTokenAlive(token string) bool {
	var request *fasthttp.Request = fasthttp.AcquireRequest()
	var response *fasthttp.Response = fasthttp.AcquireResponse()
	var CSRF string = generateRandomCSRF()

	defer fasthttp.ReleaseRequest(request)
	defer fasthttp.ReleaseResponse(response)

	request.Header.SetMethod("GET")
	request.SetRequestURI("https://api.twitter.com/1.1/account/settings.json")

	request.Header.Set("Connection", "close")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0")

	request.Header.Set("Authorization", "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA")
	request.Header.Set("X-CSRF-Token", CSRF)
	request.Header.Set("Cookie", "auth_token=" + token + "; ct0=" + CSRF)

	fasthttp.Do(request, response)

	return response.StatusCode() == 200
}

func getUsernameID(username, token string) (string, int) {
	var request *fasthttp.Request = fasthttp.AcquireRequest()
	var response *fasthttp.Response = fasthttp.AcquireResponse()	
	var CSRF string = generateRandomCSRF()

	defer fasthttp.ReleaseRequest(request)
	defer fasthttp.ReleaseResponse(response)

	request.Header.SetMethod("GET")
	request.SetRequestURI("https://api.twitter.com/1.1/users/show.json?screen_name=" + username)

	request.Header.Set("Connection", "close")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0")

	request.Header.Set("Authorization", "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA")
	request.Header.Set("X-CSRF-Token", CSRF)
	request.Header.Set("Cookie", "auth_token=" + token + "; ct0=" + CSRF)

	fasthttp.Do(request, response)

	return fastjson.GetString(response.Body(), "id_str"), response.StatusCode()
}

func uploadImage(imageFilePath, token string) (string, int) {
	var request *fasthttp.Request = fasthttp.AcquireRequest()
	var response *fasthttp.Response = fasthttp.AcquireResponse()
	var contentType, buffer = loadImage(imageFilePath)
	var CSRF string = generateRandomCSRF()

	defer fasthttp.ReleaseRequest(request)
	defer fasthttp.ReleaseResponse(response)

	request.Header.SetMethod("POST")
	request.SetRequestURI("https://upload.twitter.com/1.1/media/upload.json")

	request.Header.Set("Content-Type", contentType)
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0")

	request.Header.Set("X-CSRF-Token", CSRF)
	request.Header.Set("Cookie", "auth_token=" + token + "; ct0=" + CSRF)
	request.Header.Set("Authorization", "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA")
	
	request.SetBody(buffer.Bytes())

	fasthttp.Do(request, response)

	return fastjson.GetString(response.Body(), "media_id_string"), response.StatusCode()
}

func getLatestTweet(usernameID, token string) (string, string, int) {
	var request *fasthttp.Request = fasthttp.AcquireRequest()
	var response *fasthttp.Response = fasthttp.AcquireResponse()
	var CSRF string = generateRandomCSRF()

	defer fasthttp.ReleaseRequest(request)
	defer fasthttp.ReleaseResponse(response)

	request.Header.SetMethod("GET")
	request.Header.SetRequestURI("https://twitter.com/i/api/graphql/V1ze5q3ijDS1VeLwLY0m7g/UserTweets?variables=%7B%22userId%22%3A%22" + usernameID + "%22%2C%22count%22%3A20%2C%22includePromotedContent%22%3Atrue%2C%22withQuickPromoteEligibilityTweetFields%22%3Atrue%2C%22withVoice%22%3Atrue%2C%22withV2Timeline%22%3Atrue%7D&features=%7B%22responsive_web_graphql_exclude_directive_enabled%22%3Atrue%2C%22verified_phone_label_enabled%22%3Afalse%2C%22creator_subscriptions_tweet_preview_api_enabled%22%3Atrue%2C%22responsive_web_graphql_timeline_navigation_enabled%22%3Atrue%2C%22responsive_web_graphql_skip_user_profile_image_extensions_enabled%22%3Afalse%2C%22c9s_tweet_anatomy_moderator_badge_enabled%22%3Atrue%2C%22tweetypie_unmention_optimization_enabled%22%3Atrue%2C%22responsive_web_edit_tweet_api_enabled%22%3Atrue%2C%22graphql_is_translatable_rweb_tweet_is_translatable_enabled%22%3Atrue%2C%22view_counts_everywhere_api_enabled%22%3Atrue%2C%22longform_notetweets_consumption_enabled%22%3Atrue%2C%22responsive_web_twitter_article_tweet_consumption_enabled%22%3Afalse%2C%22tweet_awards_web_tipping_enabled%22%3Afalse%2C%22freedom_of_speech_not_reach_fetch_enabled%22%3Atrue%2C%22standardized_nudges_misinfo%22%3Atrue%2C%22tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled%22%3Atrue%2C%22rweb_video_timestamps_enabled%22%3Atrue%2C%22longform_notetweets_rich_text_read_enabled%22%3Atrue%2C%22longform_notetweets_inline_media_enabled%22%3Atrue%2C%22responsive_web_media_download_video_enabled%22%3Afalse%2C%22responsive_web_enhance_cards_enabled%22%3Afalse%7D")
	
	request.Header.Set("Connection", "close")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0")

	request.Header.Set("Authorization", "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA")
	request.Header.Set("X-CSRF-Token", CSRF)
	request.Header.Set("Cookie", "auth_token=" + token + "; ct0=" + CSRF)

	fasthttp.Do(request, response)
	var tweetID, tweetDate string = parseTweet(string(response.Body()))

	return tweetID, tweetDate, response.StatusCode()
}

func sendReply(tweeetID, token, message string, bot *Bot) (string, int) {
	var request *fasthttp.Request = fasthttp.AcquireRequest()
	var response *fasthttp.Response = fasthttp.AcquireResponse()	
	var CSRF string  = generateRandomCSRF()

	defer fasthttp.ReleaseRequest(request)
	defer fasthttp.ReleaseResponse(response)

	request.Header.SetMethod("POST")
	request.SetRequestURI("https://twitter.com/i/api/graphql/bDE2rBtZb3uyrczSZ_pI9g/CreateTweet")

	request.Header.Set("Connection", "close")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0")

	request.Header.Set("Authorization", "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA")
	request.Header.Set("X-CSRF-Token", CSRF)	
	request.Header.Set("Cookie", "auth_token=" + token + "; ct0=" + CSRF)
	
	var payload string = `{"variables":{"tweet_text":"` + message + `","reply":{"in_reply_to_tweet_id":"` + tweeetID + `"},"media":{"media_entities":[`
	if bot.ImageID == "" {
		payload += `],"possibly_sensitive":false},"semantic_annotation_ids":[]},"features":{"c9s_tweet_anatomy_moderator_badge_enabled":true,"tweetypie_unmention_optimization_enabled":true,"responsive_web_edit_tweet_api_enabled":true,"graphql_is_translatable_rweb_tweet_is_translatable_enabled":true,"view_counts_everywhere_api_enabled":true,"longform_notetweets_consumption_enabled":true,"responsive_web_twitter_article_tweet_consumption_enabled":false,"tweet_awards_web_tipping_enabled":false,"longform_notetweets_rich_text_read_enabled":true,"longform_notetweets_inline_media_enabled":true,"rweb_video_timestamps_enabled":true,"responsive_web_graphql_exclude_directive_enabled":true,"verified_phone_label_enabled":false,"freedom_of_speech_not_reach_fetch_enabled":true,"standardized_nudges_misinfo":true,"tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled":true,"responsive_web_media_download_video_enabled":false,"responsive_web_graphql_skip_user_profile_image_extensions_enabled":false,"responsive_web_graphql_timeline_navigation_enabled":true,"responsive_web_enhance_cards_enabled":false}}`
	} else {
		payload += `{"media_id":"` + bot.ImageID + `","tagged_users":[]}],"possibly_sensitive":false},"semantic_annotation_ids":[]},"features":{"c9s_tweet_anatomy_moderator_badge_enabled":true,"tweetypie_unmention_optimization_enabled":true,"responsive_web_edit_tweet_api_enabled":true,"graphql_is_translatable_rweb_tweet_is_translatable_enabled":true,"view_counts_everywhere_api_enabled":true,"longform_notetweets_consumption_enabled":true,"responsive_web_twitter_article_tweet_consumption_enabled":false,"tweet_awards_web_tipping_enabled":false,"longform_notetweets_rich_text_read_enabled":true,"longform_notetweets_inline_media_enabled":true,"rweb_video_timestamps_enabled":true,"responsive_web_graphql_exclude_directive_enabled":true,"verified_phone_label_enabled":false,"freedom_of_speech_not_reach_fetch_enabled":true,"standardized_nudges_misinfo":true,"tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled":true,"responsive_web_media_download_video_enabled":false,"responsive_web_graphql_skip_user_profile_image_extensions_enabled":false,"responsive_web_graphql_timeline_navigation_enabled":true,"responsive_web_enhance_cards_enabled":false},"queryId":"bDE2rBtZb3uyrczSZ_pI9g"}`
	}
	request.SetBodyString(payload)

	fasthttp.Do(request, response)
	var body string = string(response.Body())

	if strings.Contains(body, "BadRequest: Invalid Media") {
		bot.setImageID()
		return sendReply(tweeetID, token, message, bot)
	}
	return fastjson.GetString([]byte(body), "data", "create_tweet", "tweet_results", "result", "rest_id"), response.StatusCode()
}

func reTweet(tweetID string, token string, message string) string {
	var request *fasthttp.Request = fasthttp.AcquireRequest()
	var response *fasthttp.Response = fasthttp.AcquireResponse()	
	var CSRF string = generateRandomCSRF()

	defer fasthttp.ReleaseRequest(request)
	defer fasthttp.ReleaseResponse(response)

	request.Header.SetMethod("POST")
	request.SetRequestURI("https://twitter.com/i/api/graphql/ojPdsZsimiJrUGLR1sjUtA/CreateRetweet")

	request.Header.Set("Connection", "close")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0")

	request.Header.Set("Authorization", "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA")
	request.Header.Set("X-CSRF-Token", CSRF)	
	request.Header.Set("Cookie", "auth_token=" + token + "; ct0=" + CSRF)
	request.SetBodyString(`{"variables":{"tweet_id":"` + tweetID + `"}}`)

	fasthttp.Do(request, response)

	return string(response.Body())
}

func likeTweet(tweetID, token string) string {
	var request *fasthttp.Request = fasthttp.AcquireRequest()
	var response *fasthttp.Response = fasthttp.AcquireResponse()	
	var CSRF string = generateRandomCSRF()

	defer fasthttp.ReleaseRequest(request)
	defer fasthttp.ReleaseResponse(response)

	request.Header.SetMethod("POST")
	request.SetRequestURI("https://twitter.com/i/api/graphql/lI07N6Otwv1PhnEgXILM7A/FavoriteTweet")

	request.Header.Set("Connection", "close")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0")

	request.Header.Set("Authorization", "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA")
	request.Header.Set("X-CSRF-Token", CSRF)	
	request.Header.Set("Cookie", "auth_token=" + token + "; ct0=" + CSRF)
	request.SetBodyString(`{"variables":{"tweet_id":"` + tweetID + `"}}`)

	fasthttp.Do(request, response)

	return string(response.Body())
}

func sendTelegramWebhook(botID, chatid, message string) {
	var request *fasthttp.Request = fasthttp.AcquireRequest()
	var response *fasthttp.Response = fasthttp.AcquireResponse()

	request.Header.SetMethod("GET")
	request.SetRequestURI("https://api.telegram.org/bot" + botID + "/sendMessage?chat_id=" + chatid + "&text=" + message)

	fasthttp.Do(request, response)
	fasthttp.ReleaseRequest(request)
}