package main

import (
	"strings"
	"github.com/bwmarrin/discordgo"
	"github.com/g4stly/navi/bot"
)

func quit(navi *bot.Bot, author *discordgo.User, argc int, argv []string) (string, error) {
	navi.Quit <- 1
	return "bai baii!", nil
}

func echo(navi *bot.Bot, author *discordgo.User, argc int, argv []string) (string, error) {
	return strings.Join(argv[1:], " "), nil
}

func reloadMods(navi *bot.Bot, author *discordgo.User, argc int, argv[]string) (string, error) {
	err := navi.ReloadModules()
	if err != nil {
		return "", err
	}
	return "**reloaded modules**", nil
}

func Init() bot.Module {
	module := bot.Module{
		Name: "navi-core",
		Commands: make(map[string]bot.Command)}
	module.Commands["quit"] = quit
	module.Commands["echo"] = echo
	module.Commands["repeat after me"] = echo
	module.Commands["reload"] = reloadMods
	return module
}
