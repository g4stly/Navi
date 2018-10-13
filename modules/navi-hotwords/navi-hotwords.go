package main

import (
	"fmt"
	"time"
	"strings"
	"math/rand"
	"github.com/bwmarrin/discordgo"
	"github.com/g4stly/navi/bot"
)

const MODULE_NAME = "navi-hotwords"

// core command type
type coreCommand struct {
	names   func() []string
	usage   func(string) string
	execute func(*bot.Bot, *discordgo.User, int, []string) (string, error)
}

func (self *coreCommand) Names() []string {
	return self.names()
}

func (self *coreCommand) Usage(nameUsed string) string {
	return self.usage(nameUsed)
}

func (self *coreCommand) Execute(navi *bot.Bot, author *discordgo.User, argc int, argv []string) (string, error) {
	return self.execute(navi, author, argc, argv)
}


func Init(navi *bot.Bot) bot.Module {
	// seed prng
	seed := time.Now().Unix()
	rand.Seed(seed)

	// initialize module
	var module bot.Module
	module.Name = MODULE_NAME
	module.HotWords = make(map[string]string)
	module.HotWords["swag"] = ".echo don't say swag"
	return module
}

func Cleanup(navi *bot.Bot) error {
	return nil
}
