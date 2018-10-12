package main

import (
	"fmt"
	"time"
	"strings"
	"math/rand"
	"github.com/bwmarrin/discordgo"
	"github.com/g4stly/navi/bot"
)

const MODULE_NAME = "navi-luck"

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

// flip command
func flipNames() []string {
	return []string{"flip"}
}

func flipUsage(nameUsed string) string {
	return fmt.Sprintf("USAGE: `.%v`\nFlip a coin.", nameUsed)
}

func flipExec(navi *bot.Bot, author *discordgo.User, argc int, argv []string) (string, error) {
	flip := "heads"
	coin := rand.Float32() * 10
	if coin <= 5 {
		flip = "tails"
	}
	return flip, nil
}

// choose command
func chooseNames() []string {
	return []string{"choose"}
}

func chooseUsage(nameUsed string) string {
	return fmt.Sprintf("USAGE: `.%v [option 1] or [option 2]`\nChoose one of the two options.", nameUsed)
}

func chooseExec(navi *bot.Bot, author *discordgo.User, argc int, argv []string) (string, error) {
	if argc < 4 {
		return chooseUsage(argv[0]), nil
	}

	perspective := map[string]string {
		"my":	"your",
		"me":	"you",
		"I":	"you",
		"i":	"you",
		"our":	"your",
		"we":	"you",
		"your":	"my",
		"you":	"navi:tm:"} // differentiating between object and subject is above my pay grade

	index := 0
	for i := 1; i < argc; i++ {
		if argv[i] == "or" {
			index = i
			continue
		}
		replacement, ok := perspective[argv[i]]
		if !ok {
			continue
		}
		argv[i] = replacement
	}
	if index < 2 { // needs to atleast be the third
		return chooseUsage(argv[0]), nil
	}

	var option1 []string
	for i := 1; i < index; i++ {
		option1 = append(option1, argv[i])
	}
	var option2 []string
	for i := index+1; i < argc; i++ {
		option2 = append(option2, argv[i])
	}

	affirmations := []string {
		" Absolutely!",
		" Desu.",
		" Without a doubt.",
		" I am 100% positive.",
		" Don't you dare ask again, the answer will be the same.",
		" I have never been more sure of anything in my entire robot life.",
		" Though tbh desu I just .flip'd a coin.",
		".. maybe? Fuck if I know.",
		" But tbh I don't have a clue",
		" But who the fuck am I? Follow your heart.",
		" You really had to ask?",
		".. Definitely."}

	result := option1
	if rand.Float32() * 10 < 5 {
		result = option2
	}
	return fmt.Sprintf("%s.%s", strings.Join(result, " "), affirmations[rand.Intn(len(affirmations))]), nil
}

func Init(navi *bot.Bot) bot.Module {
	// seed prng
	seed := time.Now().Unix()
	rand.Seed(seed)

	// initialize module
	var module bot.Module
	module.Name = MODULE_NAME
	module.Commands = []bot.Command{
		&coreCommand{ // flip
			names:   flipNames,
			usage:   flipUsage,
			execute: flipExec},
		&coreCommand{ // choose
			names:   chooseNames,
			usage:   chooseUsage,
			execute: chooseExec}}
	return module
}

func Cleanup(navi *bot.Bot) error {
	return nil
}
