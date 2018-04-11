package main

import (
	"os"
	"github.com/g4stly/navi/bot"
	"github.com/g4stly/navi/common"
)

func main() {
	// initialize bot
	botToken := common.Config["bot-token"].(string)
	navi, err := bot.New(botToken)
	if err != nil {
		common.Fatal("failed to initialize bot: %v", err)
	}
	os.Exit(navi.Connect())
}
