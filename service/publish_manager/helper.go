package publishmanager

import (
	"math/rand"
	"time"
)

const (
	minCycleTime = 60   // 1 hours
	maxCycleTime = 1668 // 27.8 hours
)

// randomTime returns a random time between minCycleTime and maxCycleTime
// in minutes.
// expected average output: 864 minutes
func randomTime() time.Duration {
	minutes := rand.Int31n(maxCycleTime-minCycleTime) + minCycleTime

	return time.Duration(minutes) * time.Minute
}

// randomNum returns a random number between 0 and 1.
// expected average output: 0.5
func randomNum() int32 {
	var min int32 = 0
	var max int32 = 10
	r := rand.Int31n(max-min) + min
	if r < 5 {
		return 1
	}
	return 0
}
