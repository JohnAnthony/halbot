package uri_title

import (
	"github.com/JohnAnthony/halbot/message"
	"regexp"
	"net/http"
	"io"
	"html"
)

var httpRe = regexp.MustCompile("https?://[^\\s]*")
var titleRe = regexp.MustCompile("<title>\\s*(?P<want>.*)\\s*</title>")

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

	buf := make([]byte, 4096)
	_, err = io.ReadFull(resp.Body, buf)

	mtype := http.DetectContentType(buf)
	if len(mtype) < 9 || mtype[0:9] != "text/html" {
		return ""
	}

	matches := titleRe.FindStringSubmatch(string(buf))
	if len(matches) < 2 {
		return ""
	}

	title := html.UnescapeString(matches[1])

	return "Site Title :: " + title
}

