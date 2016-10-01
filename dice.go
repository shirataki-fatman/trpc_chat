package main

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/seehuhn/mt19937"
)

func getDice() *rand.Rand {
	r := rand.New(mt19937.New())
	r.Seed(time.Now().UnixNano())

	return r
}

func RollDice(diceNum, diceMin, diceMax int) []string {
	r := getDice()

	var result []string
	for i := 0; i < diceNum; i++ {
		d := r.Intn(diceMax-diceMin) + diceMin
		s := strconv.Itoa(d)
		result = append(result, s)
	}

	return result
}
