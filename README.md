# Navi  
## A discord bot written in Go
Navi is a discord bot with a focus on modular command development. I want to be able to dynamically at runtime load and unload modules that contain the commands the bot can run. I guess we'll find out together if this works out.  
Make sure if you're to run it to have a json file with a key `bot-token` that is your bot's token. You can specify the name/location of the config file with the `-f` flag; it defaults to `config.json` in your cwd.
