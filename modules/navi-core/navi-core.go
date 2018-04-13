package main

import (
	"fmt"
	"bytes"
	"strings"
	"github.com/bwmarrin/discordgo"
	"github.com/g4stly/navi/bot"
)

const MODULE_NAME = "navi-core"

// core command type
type coreCommand struct {
	names	func()[]string
	usage	func(string)string
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

// echo command
func echoNames() []string {
	return []string{"echo"}
}

func echoUsage(nameUsed string) string {
	return fmt.Sprintf("USAGE: `.%v ...`\nNavi will repeat what you say.", nameUsed)
}

func echoExec(navi *bot.Bot, author *discordgo.User, argc int, argv []string) (string, error) {
	return strings.Join(argv[1:], " "), nil
}

// list command
func listNames() []string {
	return []string{"list","ls","commands"}
}

func listUsage(nameUsed string) string {
	return fmt.Sprintf("USAGE: `.%v`\nList available commands.", nameUsed)
}

func listExec(navi *bot.Bot, author *discordgo.User, argc int, argv []string) (string, error) {
	var response bytes.Buffer
	for cmdName := range navi.Commands {
		_, err := response.WriteString(fmt.Sprintf(".%v\n", cmdName))
		if err != nil { return "", err}
	}
	return response.String(), nil
}

// usage command
func usageNames() []string {
	return []string{"usage","help"}
}

func usageUsage(nameUsed string) string {
	return fmt.Sprintf("USAGE: `.%v [command name]`\nWill print the usage for a given command.\nSquare brackets denote mandatory arguments, parenthesis denote optional arguments.\n", nameUsed)
}

func usageExec(navi *bot.Bot, author *discordgo.User, argc int, argv []string) (string, error) {
	if argc < 2 {
		return usageUsage(argv[0]), nil
	}
	return navi.Commands[argv[1]].Usage(argv[1]), nil
}

// combo command
func comboNames() []string {
	return []string{"combo"}
}

func comboUsage(nameUsed string) string {
	return fmt.Sprintf("USAGE: `.%v`\nToggle combo mode.", nameUsed)
}

func comboExec(navi *bot.Bot, author *discordgo.User, argc int, argv []string) (string, error) {
	mode := "disabled"
	navi.Combo = !navi.Combo
	if navi.Combo { mode = "enabled" }
	return fmt.Sprintf("Combo mode has been %v.", mode), nil
}

// quit command
func quitNames() []string {
	return []string{"quit","exit","restart","goodbye","fuckoff"}
}

func quitUsage(nameUsed string) string {
	return fmt.Sprintf("USAGE: `.%v`\nStarts a graceful shutdown.", nameUsed)
}

func quitExec(navi *bot.Bot, author *discordgo.User, argc int, argv []string) (string, error) {
	navi.Quit <- 1
	return "baibai~~", nil
}

func Init() bot.Module {
	var module bot.Module
	module.Name = MODULE_NAME
	module.Commands = []bot.Command{
		&coreCommand{	// echo
			names:echoNames,
			usage:echoUsage,
			execute:echoExec},
		&coreCommand{	// list
			names:listNames,
			usage:listUsage,
			execute:listExec},
		&coreCommand{	// usage
			names:usageNames,
			usage:usageUsage,
			execute:usageExec},
		&coreCommand{	// combo
			names:comboNames,
			usage:comboUsage,
			execute:comboExec},
		&coreCommand{	// quit
			names:quitNames,
			usage:quitUsage,
			execute:quitExec}}
	return module
}
