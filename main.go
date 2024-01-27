package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	fmt.Printf("[*] Twitter Auto-Reply \n\n")
	var bot *Bot = loadConfig()

	fmt.Printf("[*] Regular Tokens: %s \n", formatNumber(int64(len(bot.RegularTokens))))
	fmt.Printf("[*] Reply Tokens: %s \n", formatNumber(int64(len(bot.ReplyTokens))))
	fmt.Printf("[*] Targets: %s \n\n", formatNumber(int64(len(bot.Targets))))

	fmt.Printf("[!] Running background check... \n\n")
	bot.GrabUsernameIds()
	bot.setImageID()

	fmt.Printf("[!] Ready, Press ENTER to start!")
	fmt.Scanln()
	fmt.Println()

	bot.RunMonitor = true
	var channel chan string = make(chan string)
	go bot.monitorAccounts(channel)

	go func() {
		for {
			for _, username := range bot.Targets {
				channel <- username
			}
		}
	}()

	go func() {
		for {
			if len(bot.ReplyTokens) < 1 || len(bot.RegularTokens) < 1 {
				log.Printf("[-] No tokens, please restart bot with atleast 1 token")
				os.Exit(0)
			}
		}
	}()

	go func() {
		for {
			time.Sleep(time.Second * time.Duration(bot.MonitorRunPeriod))
			bot.RunMonitor = false
			time.Sleep(time.Second * time.Duration(bot.MonitorSleepDelay))
			bot.RunMonitor = true
		}
	}()

	go func() {
		if bot.ImageReply {
			for {
				time.Sleep(time.Second * 6000)
				bot.setImageID()
			}
		}
	}()

	select {}
}