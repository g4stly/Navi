package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/g4stly/navi/common"
	"plugin"
	"time"
)

// types
type Bot struct {
	session        *discordgo.Session
	ID             string
	Combo          bool
	Commands       map[string]Command
	Database       database // (defined in database.go)
	Modules        []Module
	EnabledModules map[string]bool
	Permissions    map[string]int
	Quit           chan int
}

//type Command func(*Bot, *discordgo.User, int, []string) (string, error)
type Command interface {
	Names() []string
	Usage(string) string
	Execute(*Bot, *discordgo.User, int, []string) (string, error)
}

type Module struct {
	Name     string
	Commands []Command
}

// variables
var ErrNotAuthorized = "**You aren't authorized to execute that command!**"

// methods for Bot type
func (self *Bot) loadModules() error {
	Modules := common.Config["modules"].([]interface{})
	for _, moduleName := range Modules {
		common.Log("loading module: %v", moduleName.(string))

		err := self.loadModule(moduleName.(string))
		if err != nil {
			return err
		}

		common.Log("done")
	}
	return nil
}

func (self *Bot) loadModule(moduleName string) error {
	moduleLocation := fmt.Sprintf("modules/%v/%v.so", moduleName, moduleName)
	moduleLibrary, err := plugin.Open(moduleLocation)
	if err != nil {
		return err
	}

	initialize, err := moduleLibrary.Lookup("Init")
	if err != nil {
		return err
	}

	module := initialize.(func() Module)()
	self.Modules = append(self.Modules, module)
	self.EnabledModules[module.Name] = true

	return nil
}

func (self *Bot) ReloadCommands() {
	for _, module := range self.Modules {
		common.Log("loading commands for module: %v", module.Name)
		if !self.EnabledModules[module.Name] {
			continue
		}
		for _, command := range module.Commands {
			for _, name := range command.Names() {
				common.Log("\t\tloading %v...", name)
				self.Commands[name] = command
			}
		}
	}
}

func (self *Bot) loadPermissions() {
	common.Log("Loading permissions from database...")
	rows, err := self.Database.query("SELECT userID, level FROM permissions;")
	if err != nil {
		common.Fatal("loadPermissions(): %v", err)
	}

	for rows.Next() {
		var userID string
		var level int
		rows.Scan(&userID, &level)
		self.Permissions[userID] = level
	}
	common.Log("done.")
}
func (self *Bot) savePermissions() {
	common.Log("Saving permissions into database...")
	commandString := "REPLACE INTO permissions (userID, level) VALUES ('?', '?');"
	for userID, level := range self.Permissions {
		_, err := self.Database.exec(commandString, userID, level)
		if err != nil {
			common.Fatal("savePermissions(): %v", err)
		}
	}
	common.Log("done.")
}

func (self *Bot) Connect() int {
	common.Log("Connecting to Discord...")

	err := self.session.Open()
	if err != nil {
		common.Log("Failed to connect to Discord: %v", err)
		return -1
	}
	defer self.session.Close()

	common.Log("Done")
	return self.waitForQuit()
}

func (self *Bot) waitForQuit() int {
	for {
		select {
		case <-self.Quit:
			common.Log("Gracefully shutting down")
			resultCode := self.onShutdown()
			common.Log("Done")
			return resultCode
		default:
			time.Sleep(500 * time.Millisecond)
			break
		}
	}
}

// main exported function, initializes and returns a Bot type
func New(token string) (*Bot, error) {
	common.Log("Initializing bot...")
	var err error

	// initialize bot
	bot := &Bot{
		Combo:          false,
		Commands:       make(map[string]Command),
		EnabledModules: make(map[string]bool),
		Permissions:    make(map[string]int),
		Quit:           make(chan int)}

	// initialize database
	err = bot.Database.startup(common.Config["database-location"].(string))

	// grab permissions from database
	bot.loadPermissions()

	// load modules
	err = bot.loadModules()
	if err != nil {
		common.Fatal("failed to load module: %v", err)
	}
	bot.ReloadCommands()

	// create discord session
	bot.session, err = discordgo.New("Bot " + token)
	if err != nil {
		return nil, common.NewError("failed to create discordgo session: %v", err)
	}

	// add callbacks (defined in callbacks.go)
	bot.session.AddHandler(bot.ready)
	bot.session.AddHandler(bot.messageCreate)

	common.Log("Done")
	return bot, err
}
