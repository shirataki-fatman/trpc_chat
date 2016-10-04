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
	Plugin  string
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
	log := t.getLog(args)
	log.Write(args.Name + "\t" + args.Message)

	t.send(args.Message, args.Channel)
	if !t.parseDiceCommand(args) {
		pluginResult := pluginManager.Exec(args)
		if pluginResult != "" {
			t.send(pluginResult, args.Channel)
		}
	}
	*result = 1
	return nil
}

func (t *Chat) getLog(args *ChatArgs) *Log {
	log, err := GetLog(args.Channel)
	if err != nil {
		panic(err)
	}

	return log
}

func (t *Chat) send(message, channel string) {
	t.melody.BroadcastFilter([]byte(message), func(q *melody.Session) bool {
		return q.Request.URL.Path == "/channel/"+channel
	})
}

func (t *Chat) parseDiceCommand(args *ChatArgs) bool {
	regex := regexp.MustCompile("^(\\d+)d(\\d+)\\s?")

	str := []byte(args.Message)
	submatch := regex.FindSubmatch(str)
	if len(submatch) <= 0 {
		return false
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

	log := t.getLog(args)
	resultString := "DiceBot\t" + matchString + " -> [" + diceResultString + "] -> " + strconv.Itoa(diceResultSum)
	log.Write(resultString)
	t.send(resultString, args.Channel)

	return true
}
