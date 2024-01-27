package main

import (
	"github.com/valyala/fastjson"
	"io/ioutil"
	"os"
)

type Bot struct {
	JSON *fastjson.Value

	Targets []string
	RegularTokens []string
	ReplyTokens []string
	Replied []string

	MonitorRunPeriod int
	MonitorSleepDelay int
	BotterGoroutines int

	RunMonitor bool
	RunBotter bool
	spammerSynchronize bool

	ReplyMessage string
	RetweetMessage string

	ImageReply bool
	ImageFilePath string
	ImageID string

	SuccessCounter int64
	ErorrCounter int64

	TelegramBotToken string
	TelegramChatID string
}

func loadConfig() *Bot {
	var file, _= os.Open("./data/config.json")
	var bytes, _ = ioutil.ReadAll(file)
	var json, _  = fastjson.Parse(string(bytes))
	return &Bot {
		ReplyTokens: openFile("./data/reply_tokens.txt"),
		RegularTokens: openFile("./data/regular_tokens.txt"),
		Targets: openFile("./data/targets.txt"),
		Replied: openFile("./data/replied.txt"),

		MonitorRunPeriod: json.GetInt("monitor_run_period"),
		MonitorSleepDelay: json.GetInt("monitor_sleep_delay"),
		BotterGoroutines: json.GetInt("spammer_goroutines"),

		ReplyMessage: string(json.GetStringBytes("bot_reply_message")),
		RetweetMessage: string(json.GetStringBytes("bot_retweet_message")),

		ImageReply: bool(json.GetBool("image", "image_reply")),
		ImageFilePath: string(json.GetStringBytes("image", "image_file_path")),
		
		TelegramBotToken: string(json.GetStringBytes("telegram_bot_token")),
		TelegramChatID: string(json.GetStringBytes("telegram_chat_id")),
	}
}