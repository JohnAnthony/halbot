package uri_title

import (
	"halbot/message"
	"regexp"
	"net/http"
	"io"
)

var httpRe = regexp.MustCompile("https?://[^ ]*")
var titleRe = regexp.MustCompile("<title>(?P<want>.*)</title>")

func Handler(m message.Message) string {
	if m.Type != "PRIVMSG" {
		return ""
	}

	url := httpRe.FindString(m.Contents)
	if url == "" {
		return ""
	}

	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	buf := make([]byte, 512)
	_, err = io.ReadFull(resp.Body, buf)

	mtype := http.DetectContentType(buf)
	if len(mtype) < 9 || mtype[0:9] != "text/html" {
		return "MIME Type: " + mtype
	}

	matches := titleRe.FindStringSubmatch(string(buf))
	if len(matches) < 2 {
		return ""
	}

	title := matches[1]
	return "Site Title: " + title
}
