package main

import (
	"bufio"
	"os"
	"sync"
)

var channels = map[string]*Log{}

func GetLog(channel string) (*Log, error) {
	log, ok := channels[channel]
	if !ok {
		var err error
		log, err = NewLog(channel)
		if err != nil {
			return nil, err
		}
		channels[channel] = log
	}

	return log, nil
}

type Log struct {
	fileLock sync.Mutex
	channel  string
}

func NewLog(channel string) (*Log, error) {
	log := Log{
		channel: channel,
	}

	_, err := os.Stat(log.filePath())
	if err != nil {
		_, createErr := os.Create(log.filePath())
		if createErr != nil {
			return nil, createErr
		}
	}

	return &log, nil
}

func (l *Log) Load() ([]string, error) {
	var logArr []string
	fp, err := os.Open(l.filePath())
	defer fp.Close()
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		log := scanner.Text()
		logArr = append(logArr, log)
	}

	return logArr, nil
}

func (l *Log) Write(message string) error {
	l.fileLock.Lock()
	defer l.fileLock.Unlock()

	fp, err := os.OpenFile(l.filePath(), os.O_APPEND|os.O_WRONLY, 0600)
	defer fp.Close()
	if err != nil {
		return err
	}
	_, err = fp.WriteString(message + "\n")
	if err != nil {
		return err
	}

	return nil
}

func (l *Log) filePath() string {
	return "./log/" + l.channel + ".log"
}
