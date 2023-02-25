package util

import (
	"time"
)

func Schedule(initDuration time.Duration, procedure func() time.Duration) {
	ticker := time.NewTicker(initDuration)
	go func(procedure func() time.Duration) {
		for {
			<-ticker.C
			nextDuration := procedure()
			ticker.Reset(nextDuration)
		}
	}(procedure)
}
