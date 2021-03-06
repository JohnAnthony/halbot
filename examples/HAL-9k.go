package main

import (
	"github.com/JohnAnthony/halbot"
	"github.com/JohnAnthony/halbot/uri_title"
)

func main() {
	myBot := halbot.NewHALBot("HAL-9k", "irc.rizon.net", 6660, "#/g/sicp")
	myBot.AddHandler(uri_title.Handler)
	myBot.Run()
}
