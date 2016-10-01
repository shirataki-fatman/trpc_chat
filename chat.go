package main

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

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

	t.send(args.Message, args.Channel)
	t.parseDiceCommand(args)
	*result = 1
	return nil
}

func (t *Chat) send(message, channel string) {
	t.melody.BroadcastFilter([]byte(message), func(q *melody.Session) bool {
		return q.Request.URL.Path == "/channel/"+channel
	})
}

func (t *Chat) parseDiceCommand(args *ChatArgs) {
	regex := regexp.MustCompile("^(\\d+)d(\\d+)\\s?")

	str := []byte(args.Message)
	submatch := regex.FindSubmatch(str)
	if len(submatch) <= 0 {
		return
	}

	var matchString string
	var diceNum, diceMax, diceMin int
	var err error

	matchString = string(submatch[0])
	diceNum, err = strconv.Atoi(string(submatch[1]))
	if err != nil {
		panic(err)
	}
	diceMax, err = strconv.Atoi(string(submatch[2]))
	if err != nil {
		panic(err)
	}
	if diceMax >= 10 {
		diceMax = diceMax - 1
		diceMin = 0
	} else {
		diceMin = 1
	}

	diceResult := RollDice(diceNum, diceMin, diceMax)
	diceResultString := strings.Join(diceResult, " + ")
	diceResultSum := 0
	for d := range diceResult {
		num, err := strconv.Atoi(diceResult[d])
		if err != nil {
			panic(err)
		}
		diceResultSum += num
	}

	resultString := matchString + " -> [" + diceResultString + "] -> " + strconv.Itoa(diceResultSum)
	t.send(resultString, args.Channel)
}
