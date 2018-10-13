package main

import (
	"fmt"
	"time"
	"math/rand"
	"strings"
	"net/url"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"github.com/g4stly/navi/bot"
	"github.com/g4stly/navi/common"
)

const MODULE_NAME = "navi-reaction"
var laffinLinks []string
var definitions map[string]definition

// definition type 
type definition struct {
	Author		string `json:"author"`	// discord user id
	Word		string `json:"key"`		// target word 
	Definition	string `json:"value"`	// definition
}

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

// define command
func defineNames() []string {
	return []string{"define", "def"}
}

func defineUsage(nameUsed string) string {
	return fmt.Sprintf("USAGE: `.%v <word> (as <new definition>)`\nView/Assign a word's custom definition.", nameUsed)
}

func defineExec(navi *bot.Bot, author *discordgo.User, argc int, argv []string) (string, error) {
	if argc < 2 {
		return defineUsage(argv[0]), nil
	}
	if argc > 2 { // define new
		if argc < 4 || argv[2] != "as" {
			return defineUsage(argv[0]), nil
		}
		def, ok := definitions[argv[1]]
		if ok && def.Author != author.ID {
			return fmt.Sprintf("**%v** is already defined as %v", def.Word, def.Definition), nil
		}
		definitions[argv[1]] = definition {
			Author: author.ID,
			Word: argv[1],
			Definition: strings.Join(argv[3:], " ")}
		return "Got it! :3", nil
	}
	// fetch definition
	def, ok := definitions[argv[1]]
	if !ok {
		return fmt.Sprintf("There doesn't seem to be a definition for %v yet..", argv[1]), nil
	}
	return fmt.Sprintf("**%v**: %v", def.Word, def.Definition), nil
}

func Init(navi *bot.Bot) bot.Module {
	// seed prng
	seed := time.Now().Unix()
	rand.Seed(seed)

	// load laffin links from database
	var err error
	laffinLinks, err = navi.Database.LoadSlice("laffin")
	if err != nil {
		common.Log("failed to load database for laffin command;\n%v\ncontinuing...", err)
		laffinLinks = make([]string, 0)
	}

	// load definition map from database
	definitions = make(map[string]definition)
	definitionSlice, err := navi.Database.LoadSlice("definitions")
	if err != nil {
		common.Error("failed to load resources for definition command;\n%v\ncontinuing...", err)
		definitionSlice = make([]string, 0)
	}
	for _, rawJson := range definitionSlice {
		var def definition
		err = json.Unmarshal([]byte(rawJson), &def)
		if err != nil {
			common.Error("failed to unmarshal definition;\n%v\ncontinuing...", err)
			continue
		}
		definitions[def.Word] = def
	}

	// initialize module
	var module bot.Module
	module.Name = MODULE_NAME
	module.Commands = []bot.Command{
		&coreCommand{ // define
			names:   defineNames,
			usage:   defineUsage,
			execute: defineExec},
		&coreCommand{ // laffin
			names:   laffinNames,
			usage:   laffinUsage,
			execute: laffinExec}}
	module.HotWords = make(map[string]string)
	return module
}

func Cleanup(navi *bot.Bot) error {
	defSlice := make([]string, 0)
	for _, def := range definitions {
		defJson, err := json.Marshal(def)
		if err != nil {
			common.Error("failed to marshal definition;\n%v\ncontinuing...", err)
			continue
		}
		defSlice = append(defSlice, string(defJson))
	}
	navi.Database.SaveSlice("definitions", defSlice)
	navi.Database.SaveSlice("laffin", laffinLinks)
	return nil
}
