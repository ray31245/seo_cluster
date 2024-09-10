package commentbot

import (
	"crypto/rand"
	"math/big"
	"time"
)

const (
	minCycleTime = 30 // 30 minutes
	maxCycleTime = 60 // 60 minutes
)

// randomTime returns a random time between minCycleTime and maxCycleTime
// in minutes.
func randomTime() time.Duration {
	nBig, err := rand.Int(rand.Reader, big.NewInt(maxCycleTime-minCycleTime))
	if err != nil {
		panic(err)
	}

	minutes := nBig.Int64() + minCycleTime

	return time.Duration(minutes) * time.Minute
}

// randomNum returns a random number between 0 and 100.
func randomNum() int32 {
	var min int64

	var max int64 = 100

	nBig, err := rand.Int(rand.Reader, big.NewInt(max-min))
	if err != nil {
		panic(err)
	}

	r := nBig.Int64()

	return int32(r)
}
