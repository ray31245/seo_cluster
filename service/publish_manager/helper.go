package publishmanager

import (
	"crypto/rand"
	"math/big"
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
	nBig, err := rand.Int(rand.Reader, big.NewInt(maxCycleTime-minCycleTime))
	if err != nil {
		panic(err)
	}

	minutes := nBig.Int64() + minCycleTime

	return time.Duration(minutes) * time.Minute
}

// randomNum returns a random number between 0 and 1.
// expected average output: 0.5
func randomNum() int32 {
	var min int64

	var max int64 = 10

	nBig, err := rand.Int(rand.Reader, big.NewInt(max-min))
	if err != nil {
		panic(err)
	}

	var chance int64 = 5

	r := nBig.Int64()
	if r < chance {
		return 1
	}

	return 0
}
