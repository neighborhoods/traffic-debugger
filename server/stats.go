package main

import (
	"os"
	"time"
)

var uuid = mustRandomHex(4)
var startTime time.Time = time.Now()

type stats struct {
	EnvVars             map[string]*string `json:"envVars"`
	UUID                string             `json:"uuid"`
	SleepDuration       time.Duration      `json:"sleepDuration"`
	SleepDurationString string             `json:"sleepDurationString"`
	StartTime           time.Time          `json:"startTime"`
	RequestTime         time.Time          `json:"requestTime"`
	RandomString        string             `json:"random"`
}

func generateStats(envVars []string) stats {
	s := stats{}

	s.EnvVars = getEnvVarMapping(envVars)
	s.UUID = uuid
	s.StartTime = startTime
	s.RequestTime = time.Now()
	s.RandomString = mustRandomHex(4)

	return s
}

func getEnvVarMapping(envVars []string) map[string]*string {
	vars := make(map[string]*string)

	for _, varToEcho := range envVars {
		value, ok := os.LookupEnv(varToEcho)
		if ok {
			vars[varToEcho] = &value
		} else {
			vars[varToEcho] = nil
		}
	}

	return vars
}
