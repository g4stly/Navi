TARGET 	= main
MODULES	= modules/navi-core/navi-core.so

all: $(TARGET) $(MODULES)

$(TARGET): common/common.go bot/bot.go bot/callbacks.go main.go
	go build main.go
%.so: bot/bot.go bot/callbacks.go
	go build -buildmode=plugin -o $@ $(@:%.so=%.go)

.PHONY: clean test
clean: 
	rm -vf $(TARGET) $(MODULES)
