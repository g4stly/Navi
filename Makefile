TARGET 	= main
BOT	= bot/bot.go bot/callbacks.go bot/database.go
CORE	= modules/navi-core/navi-core.so
LUCK	= modules/navi-luck/navi-luck.so
<<<<<<< HEAD
REACTION= modules/navi-reaction/navi-reaction.so
MODULES	= $(CORE) $(LUCK) $(REACTION)
=======
REACT	= modules/navi-reaction/navi-reaction.so
HOTWORDS = modules/navi-hotwords/navi-hotwords.so
MODULES	= $(CORE) $(LUCK) $(REACT) $(HOTWORDS)
>>>>>>> f99e2f440640d67d243705067a762cef4d6afef9

all: $(TARGET) $(MODULES)

$(TARGET): common/common.go $(BOT)
	go build main.go
$(CORE): $(CORE:%.so=%.go) $(BOT)
	go build -buildmode=plugin -o $@ $(@:%.so=%.go)
$(LUCK): $(LUCK:%.so=%.go) $(BOT)
	go build -buildmode=plugin -o $@ $(@:%.so=%.go)
$(REACTION): $(REACTION:%.so=%.go) $(BOT)
	go build -buildmode=plugin -o $@ $(@:%.so=%.go)
$(HOTWORDS): $(HOTWORDS:%.so=%.go) $(BOT)
	go build -buildmode=plugin -o $@ $(@:%.so=%.go)

.PHONY: clean test
clean: 
	rm -vf $(TARGET) $(MODULES)
