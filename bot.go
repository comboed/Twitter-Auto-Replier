package main

import (
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
	"fmt"
	"log"
	"os"
)

func (bot *Bot) monitorAccounts(channel chan string) {
	for {
		if bot.RunMonitor {
			for _, target := range bot.Targets {
				var username []string = strings.Split(target, ":")
				var token string = bot.RegularTokens[rand.Intn(len(bot.RegularTokens))]				
				var tweetID, tweetDate, statusCode = getLatestTweet(username[1], token)

				var currentTime time.Time = time.Now()
				var currentDay int = currentTime.Day()
				
				if strings.Contains(tweetDate, currentTime.Month().String()[:3]) {
					var tweetDay, _ = strconv.Atoi(tweetDate[8:][:2])
					if (tweetDay >= currentDay - 1 || tweetDay <= currentDay + 1) && !containsString(bot.Replied, tweetID) {
						log.Printf("Found new tweet! [Username: %q] [Tweet ID: %q] \n", username[0], tweetID)
						bot.replyToTweet(tweetID, tweetDate, username[0])
					}
				} else if statusCode == 404 {
					bot.HandleInvalidUsername(username[0], username[1])
				} else if statusCode != 200 {
					bot.HandleInvalidToken(token)
				}
			}
		}
	}
}

func (bot *Bot) replyToTweet(tweetID, tweetDate, username string) {
	for {
		var token string = bot.ReplyTokens[rand.Intn(len(bot.ReplyTokens))]
		var replyTweetID, statusCode = sendReply(tweetID, token, bot.ReplyMessage, bot)

		if statusCode == 200 {
			log.Printf("Successfully replied to tweet! [Tweet ID: %s] \n", replyTweetID)
			bot.Replied = append(bot.Replied, tweetID)

			bot.RunBotter = true
			bot.runTweetBotter(replyTweetID, tweetDate, username)
			writeFile(tweetID, "./data/replied.txt")
			break
		} else if statusCode == 400 {
			log.Printf("ERROR: Tweet ID [%s] is invalid/removed \n", tweetID)
			bot.Replied = removeString(bot.Replied, tweetID)
			break
		} else if statusCode != 429 {
			bot.HandleInvalidToken(token)
		}
	}
}

func (bot *Bot) runTweetBotter(tweetID, tweetDate, username string) {
	var waitGroup sync.WaitGroup
	var channel chan string = createChannel(bot.ReplyTokens)
	for i := 0; i < bot.BotterGoroutines && bot.RunBotter; i++ {
		waitGroup.Add(1)
		go func() {
			bot.tweetSpammer(tweetID, channel)
			waitGroup.Done()
		}()
	}
	waitGroup.Wait()
	bot.RunBotter = false
	log.Printf("Finished running tweet botter! [Success: %s | Error's: %s] \n\n", formatNumber(bot.SuccessCounter), formatNumber(bot.ErorrCounter))
	sendTelegramWebhook(bot.TelegramBotToken, bot.TelegramChatID, "Successfuly%20botted%20tweet!%0ATweet%20ID:%20" + tweetID + "%0ATweet%20Date:%20" + tweetDate + "%0ATwitter%20Handle:%20" + username)
	bot.SuccessCounter, bot.ErorrCounter = 0, 0
}

func (bot *Bot) tweetSpammer(tweetID string, channel chan string) {
	for len(channel) != 0 && bot.RunBotter {
		var token string = <- channel
		var retweetResponse string = reTweet(tweetID, token, bot.ReplyMessage)
		var likeResponse string = likeTweet(tweetID, token)

		if strings.Contains(retweetResponse, `"full_text":"`) && strings.Contains(likeResponse, `tweet":"Done"`) {
			bot.SuccessCounter++

		} else if strings.Contains(retweetResponse, `No status found`) || strings.Contains(likeResponse, `Missing: Tweet record`) && bot.RunBotter {
			bot.RunBotter = false
			fmt.Printf("[-] Tweet ID [%s] is invalid/removed \n", tweetID)

		} else if !bot.spammerSynchronize {
			bot.spammerSynchronize = true
			bot.HandleInvalidToken(token)
			bot.ErorrCounter++
			bot.spammerSynchronize = false
		}
	}
}

func (bot *Bot) GrabUsernameIds() {
	var index int
	for index < len(bot.Targets) {
		if !strings.Contains(bot.Targets[index], ":") {
			var token string = bot.RegularTokens[rand.Intn(len(bot.RegularTokens))]
			var usernameID, statusCode = getUsernameID(bot.Targets[index], token)

			if statusCode == 200 {
				bot.Targets[index] = bot.Targets[index] + ":" + usernameID
				index++
			} else if statusCode == 400 {
				bot.Targets = removeString(bot.Targets, bot.Targets[index])
				index++
			} else if statusCode != 429 {
				bot.HandleInvalidToken(token)
				continue
			}
			createFile(bot.Targets, "./data/targets.txt")
		} else {
			index++
		}
	}
}

func (bot *Bot) setImageID() {
	if bot.ImageReply {
		for {
			var token string = bot.RegularTokens[rand.Intn(len(bot.RegularTokens))]
			var imageID, statusCode = uploadImage(bot.ImageFilePath, token)
			if statusCode == 200 {
				bot.ImageID = imageID
				break
			} else if statusCode == 400 {
				fmt.Printf("[-] Failed to upload image to twitter \n")
				os.Exit(0)
			} else if statusCode != 429 {
				bot.HandleInvalidToken(token)
			}
		}
	}
}

func (bot *Bot) HandleInvalidToken(token string) {
	if !isTokenAlive(token) {
		fmt.Printf("[-] Removed invalid token [%s] \n", token)
		bot.ReplyTokens = removeString(bot.ReplyTokens, token)
		bot.RegularTokens = removeString(bot.RegularTokens, token)
		
		createFile(bot.ReplyTokens, "./data/reply_tokens.txt")
		createFile(bot.RegularTokens, "./data/regular_tokens.txt")
	
		writeFile(token, "./data/dead_tokens.txt")
	}
}

func (bot *Bot) HandleInvalidUsername(username, previousUsernameID string) {
	bot.Targets = removeString(bot.Targets, username + ":" + previousUsernameID)
	var token string = bot.RegularTokens[rand.Intn(len(bot.RegularTokens))]
	var usernameID, statusCode = getUsernameID(username, token)
	for {		
		if statusCode == 200 {
			bot.Targets = append(bot.Targets, username + ":" + usernameID)
			createFile(bot.Targets, "./data/usernames.txt")
			break
		} else if statusCode == 400 {
			log.Printf("Username %q has been suspended \n", username[0])
			break
		}  else {
			bot.HandleInvalidToken(token)
		}
	}
}