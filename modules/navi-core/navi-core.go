package main

import (
	"bytes"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/g4stly/navi/bot"
	"strconv"
	"strings"
)

const MODULE_NAME = "navi-core"

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
	return []string{"list", "ls", "commands"}
}

func listUsage(nameUsed string) string {
	return fmt.Sprintf("USAGE: `.%v`\nList available commands.", nameUsed)
}

func listExec(navi *bot.Bot, author *discordgo.User, argc int, argv []string) (string, error) {
	var response bytes.Buffer
	// loop over modules first
	for modIndex := range navi.Modules {
		// if it's disabled, skip it
		if !navi.EnabledModules[navi.Modules[modIndex].Name] {
			continue
		}
		_, err := response.WriteString(navi.Modules[modIndex].Name)
		if err != nil {
			return "", err
		}
		// now loop over commands in module
		for _, cmd := range navi.Modules[modIndex].Commands {
			// print all aliases of a command
			first := true
			names := cmd.Names()
			for _, cmdName := range names {
				if first { // decided to print comma or indent
					_, err = response.WriteString(fmt.Sprintf("\n\t`.%v`", cmdName))
					if err != nil {
						return "", err
					}
					first = false
					continue
				}
				_, err = response.WriteString(fmt.Sprintf(" `.%v`", cmdName))
				if err != nil {
					return "", err
				}
			}
		}
		_, err = response.WriteString("\n")
		if err != nil {
			return "", err
		}
	}
	return response.String(), nil
}

// usage command
func usageNames() []string {
	return []string{"usage", "help"}
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

// modToggle command
func modToggleNames() []string {
	return []string{"modToggle"}
}

func modToggleUsage(nameUsed string) string {
	return fmt.Sprintf("USAGE: `.%v [module name]`\nWill enable/disable a given module.", nameUsed)
}

func modToggleExec(navi *bot.Bot, author *discordgo.User, argc int, argv []string) (string, error) {
	if navi.Permissions[author.ID] < 1 {
		return bot.ErrNotAuthorized, nil
	}
	if argc < 2 {
		return modToggleUsage(argv[0]), nil
	}
	if argv[1] == MODULE_NAME {
		return "You can't disable the core module baka.", nil
	}
	if _, ok := navi.EnabledModules[argv[1]]; !ok {
		return fmt.Sprintf("%v: no such module. baka.", argv[1]), nil
	}
	mode := "disabled"
	navi.EnabledModules[argv[1]] = !navi.EnabledModules[argv[1]]
	if navi.EnabledModules[argv[1]] {
		mode = "enabled"
	}
	navi.ReloadCommands()
	return fmt.Sprintf("%v has been %v.", argv[1], mode), nil
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
	if navi.Combo {
		mode = "enabled"
	}
	return fmt.Sprintf("Combo mode has been %v.", mode), nil
}

// authorize command
func authorizeNames() []string {
	return []string{"authorize"}
}

func authorizeUsage(nameUsed string) string {
	return fmt.Sprintf("USAGE: `.%v (@mention (authority level))`\nManage permission levels.", nameUsed)
}

func isMention(mention string) bool {
	// todo: check if the rest of the string
	// fits in a reg exp
	if mention[0] == '@' {
		return true
	}
	return false
}

func mentionToID(mention string) string {
	return mention[2 : len(mention)-1]
}

func authorizeExec(navi *bot.Bot, author *discordgo.User, argc int, argv []string) (string, error) {
	if argc < 2 {
		return fmt.Sprintf("You have a permission level of %v.", navi.Permissions[author.ID]), nil
	} else if argc < 3 {
		if !isMention(argv[1]) {
			return authorizeUsage(argv[0]), nil
		}
		return fmt.Sprintf("%v has a permission level of %v.", argv[1], navi.Permissions[mentionToID(argv[1])]), nil
	}
	// protection cases
	if navi.Permissions[author.ID] < 2 {
		return bot.ErrNotAuthorized, nil
	}
	if navi.Permissions[author.ID] <= navi.Permissions[mentionToID(argv[1])] {
		return fmt.Sprintf("**You cannot modify the permission level of someone with greater or equal permissions.**"), nil
	}
	newLevel, err := strconv.Atoi(argv[2])
	if err != nil {
		return "", err
	}
	if newLevel < 0 {
		newLevel = 0
	}
	if newLevel > navi.Permissions[author.ID] {
		return fmt.Sprintf("**You cannot set a users permission level higher than your own.**"), nil
	}

	navi.Permissions[mentionToID(argv[1])] = newLevel
	return fmt.Sprintf("%v now has level %v permissions.", argv[1], navi.Permissions[mentionToID(argv[1])]), nil
}

// quit command
func quitNames() []string {
	return []string{"quit", "exit", "restart", "fuckoff"}
}

func quitUsage(nameUsed string) string {
	return fmt.Sprintf("USAGE: `.%v`\nStarts a graceful shutdown.", nameUsed)
}

func quitExec(navi *bot.Bot, author *discordgo.User, argc int, argv []string) (string, error) {
	if navi.Permissions[author.ID] < 3 {
		return bot.ErrNotAuthorized, nil
	}
	navi.Quit <- 1
	return "baibai~~", nil
}

func Init(navi *bot.Bot) bot.Module {
	var module bot.Module
	module.Name = MODULE_NAME
	module.Commands = []bot.Command{
		&coreCommand{ // echo
			names:   echoNames,
			usage:   echoUsage,
			execute: echoExec},
		&coreCommand{ // list
			names:   listNames,
			usage:   listUsage,
			execute: listExec},
		&coreCommand{ // usage
			names:   usageNames,
			usage:   usageUsage,
			execute: usageExec},
		&coreCommand{ // modToggle
			names:   modToggleNames,
			usage:   modToggleUsage,
			execute: modToggleExec},
		&coreCommand{ // authorize
			names:   authorizeNames,
			usage:   authorizeUsage,
			execute: authorizeExec},
		&coreCommand{ // combo
			names:   comboNames,
			usage:   comboUsage,
			execute: comboExec},
		&coreCommand{ // quit
			names:   quitNames,
			usage:   quitUsage,
			execute: quitExec}}
	module.HotWords = make(map[string]string)
	module.HotWords["nya"] = ".echo owo"
	module.HotWords["nyaa"] = ".echo ow0"
	module.HotWords["nyaaa"] = ".echo 0wo"
	module.HotWords["nyaaaa"] = ".echo 0w0"
	module.HotWords["nyaaaaa"] = ".echo 0w0"
	return module
}

func Cleanup(navi *bot.Bot) error {
	return nil
}
