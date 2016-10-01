package main

import (
	"net/http"

	"github.com/olahol/melody"
)

type ChatArgs struct {
	Message string
	Name    string
	Channel string
}

type Chat struct {
	melody *melody.Melody
}

type Result int

func newChat(melody *melody.Melody) *Chat {
	instance := Chat{
		melody: melody,
	}

	return &instance
}

func (t *Chat) SendMessage(r *http.Request, args *ChatArgs, result *Result) error {
	log, err := GetLog(args.Channel)
	if err != nil {
		panic(err)
	}
	log.Write(args.Name + "\t" + args.Message)

	t.melody.BroadcastFilter([]byte(args.Message), func(q *melody.Session) bool {
		return q.Request.URL.Path == "/channel/"+args.Channel
	})
	*result = 1
	return nil
}
