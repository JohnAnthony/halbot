package halbot

import (
	"log"
	"net"
	"strconv"
	"bufio"
	"fmt"
	"net/textproto"
	"container/list"
	"github.com/JohnAnthony/halbot/message"
)

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

type HALBot struct {
	nick string
	user string
	port int
	server string
	channel string
	conn net.Conn
	handlers list.List
}

func NewHALBot(nick string, uri string, port int, channel string) (*HALBot) {
	return &HALBot {
		nick: nick,
		user: "HALBot",
		port: port,
		server: uri,
		channel: channel,
	}
}

func (hb *HALBot) AddHandler(handler func(message.Message) string) {
	hb.handlers.PushBack(handler)
}

func (hb *HALBot) SendRaw(in string) {
	sz := min(510, len(in))
	out := in[:sz] + "\r\n"
	fmt.Printf(">> %s", out)
	hb.conn.Write([]byte(out))
}

func (hb *HALBot) connect() {
	conn, err := net.Dial("tcp", hb.server + ":" + strconv.Itoa(hb.port))
	if err != nil {
		log.Fatal("Unable to connect to IRC server: ", err)
	}
	log.Printf("Connected to IRC server %s (%s)\n", hb.server, conn.RemoteAddr())
	hb.conn = conn
}

func (hb *HALBot) disconnect() {
	hb.conn.Close()
}

func (hb *HALBot) SendToChannel(in string) {
	msg := strings.Replace(in, "\r", "", -1)
	msg = strings.Replace(msg, "\n", "", -1)
	msg = ":" + hb.nick + " PRIVMSG " + hb.channel + " :" + msg
	hb.SendRaw(msg)
}

func (hb *HALBot) Run() {
	hb.connect()
	defer hb.disconnect()

	hb.SendRaw("USER " + hb.nick + " 8 * :" + hb.nick)
	hb.SendRaw("NICK " + hb.nick)

	reader := bufio.NewReader(hb.conn)
	tp := textproto.NewReader(reader)
	for {
		// Get and display input
		line, err := tp.ReadLine()
		if err != nil {
			break // break loop on errors    
		}
		fmt.Printf("<< %s\n", line)

		// Special Case for PING -> PONG
		if (line[0:4] == "PING") {
			hb.SendRaw("PONG" + line[4:])
			continue
		}

		// Get out message
		msg := message.LineToMessage(line)

		// Special Cases
		// Join the channel after end of MOTD "376"
		if (msg.Type == "376") {
			hb.SendRaw("JOIN " + hb.channel)
			continue
		}
		// Ignore private messages
		if (msg.To != hb.channel) {
			continue
		}

		// User-defined handlers
		for e := hb.handlers.Front() ; e != nil ; e = e.Next() {
			response := e.Value.(func(message.Message) string)(msg)
			if response != "" {
				hb.SendToChannel(response)
			}
		}
	}
}
