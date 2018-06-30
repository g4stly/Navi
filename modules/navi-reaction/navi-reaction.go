package main

import (
	"fmt"
	"time"
	"math/rand"
	"net/url"
	"github.com/bwmarrin/discordgo"
	"github.com/g4stly/navi/bot"
	"github.com/g4stly/navi/common"
)

const MODULE_NAME = "navi-reaction"
var laffinLinks []string

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

// laffin command
func laffinNames() []string {
	return []string{"laffin"}
}

func laffinUsage(nameUsed string) string {
	return fmt.Sprintf("USAGE: `.%v (new <link-to-image>)`\no im laffin", nameUsed)
}

func laffinExec(navi *bot.Bot, author *discordgo.User, argc int, argv []string) (string, error) {
	if argc < 2 {
		return laffinLinks[rand.Intn(len(laffinLinks))], nil
	}
	if argv[1] != "new" || (argc < 3) {
		return laffinUsage(argv[0]), nil
	}

	newLink, err := url.Parse(argv[2])
	if err != nil {
		common.Error("url.Parse(): %v", err)
		return "That doesn't look like a url anon..", nil
	}
	if newLink.Host == "" || newLink.Scheme == "" || newLink.Path == "" {
		return "That doesn't look like a url anon..", nil
	}
	laffinLinks = append(laffinLinks, newLink.String())
	return fmt.Sprintf("Added %v to the laffin list.", newLink.String()), nil
}

func Init(navi *bot.Bot) bot.Module {
	// seed prng
	seed := time.Now().Unix()
	rand.Seed(seed)

	// load from database
	var err error
	laffinLinks, err = navi.Database.LoadSlice("laffin")
	if err != nil {
		common.Log("failed to load database for laffin command;\n%v\ncontinuing...")
		laffinLinks = make([]string, 0)
	}

	// initialize module
	var module bot.Module
	module.Name = MODULE_NAME
	module.Commands = []bot.Command{
		&coreCommand{ // laffin
			names:   laffinNames,
			usage:   laffinUsage,
			execute: laffinExec}}
	return module
}

func Cleanup(navi *bot.Bot) error {
	navi.Database.SaveSlice("laffin", laffinLinks)
	return nil
}
