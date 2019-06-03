package main

import (
	crand "crypto/rand"
	"encoding/hex"
	"math/rand"
	"time"
)

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := crand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func mustRandomHex(n int) string {
	s, err := randomHex(n)
	if err != nil {
		panic(err)
	}

	return s
}

func randomSleep(min, max time.Duration) time.Duration {
	if min >= max {
		time.Sleep(min)
		return min
	}

	var sleepRange = max - min
	var r = rand.Int63n(sleepRange.Nanoseconds())

	var sleepTime = time.Duration(r) + min
	time.Sleep(sleepTime)

	return sleepTime
}
