package message

import (
	"regexp"
)

type Message struct {
	From string
	Type string
	To string
	Contents string
}

var re = regexp.MustCompile(":(?P<from>[^ ]+) (?P<type>[^ ]+) (?P<to>[^ ]+) (?P<rest>.*)")

func LineToMessage(in string) Message {
	var m Message
	
	if in[0] != ':' {
		return Message {
			Type: "SPECIAL",
			Contents: in,
		}
	}

	m.From = re.ReplaceAllString(in, "${from}")
	m.Type = re.ReplaceAllString(in, "${type}")
	m.To = re.ReplaceAllString(in, "${to}")
	m.Contents = re.ReplaceAllString(in, "${rest}")

	return m
}
