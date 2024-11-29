package publishmanager

import (
	"context"
	"crypto/rand"
	"math/big"
	"sort"
	"time"
)

const (
	minCycleTime = 60   // 1 hours
	maxCycleTime = 1668 // 27.8 hours

	CycleTime12hours = 720  // 12 hours
	CycleTime24hours = 1440 // 24 hours
)

// randomTime returns a random time between minCycleTime and maxCycleTime
// in minutes.
// expected average output: 864 minutes
func randomTime() time.Duration {
	return randomMinute(maxCycleTime, minCycleTime)
}

func randomMinute(max, min int) time.Duration {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	if err != nil {
		panic(err)
	}

	minutes := nBig.Int64() + int64(min)

	return time.Duration(minutes) * time.Minute
}

func randomTime12hours() time.Duration {
	return randomMinute(CycleTime12hours, 0)
}

func randomTime24hours() time.Duration {
	return randomMinute(CycleTime24hours, 0)
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

func computeTimePointArray() []time.Time {
	// find next 12 or 24 clock time
	currentTime := time.Now()
	currentHour := currentTime.Hour()

	var nextBase time.Time
	if currentHour < 12 {
		nextBase = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 12, 0, 0, 0, currentTime.Location())
	} else {
		nextBase = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location()).AddDate(0, 0, 1)
	}

	timeArr := []time.Time{}

	for range 2 {
		timeArr = append(timeArr, nextBase.Add(randomTime12hours()))
	}

	for range 2 {
		timeArr = append(timeArr, nextBase.Add(12*time.Hour).Add(randomTime12hours()))
	}

	if randomNum() == 1 {
		timeArr = append(timeArr, nextBase.Add(randomTime24hours()))
	}

	sort.Slice(timeArr, func(i, j int) bool {
		return timeArr[i].Before(timeArr[j])
	})

	return timeArr
}

func timeArrSchedule(ctx context.Context, timeArr []time.Time, f func()) {
	sort.Slice(timeArr, func(i, j int) bool {
		return timeArr[i].Before(timeArr[j])
	})

	for _, t := range timeArr {
		timeToNext := time.Until(t)
		select {
		case <-ctx.Done():
			// Exit the loop if the context is cancelled
			return
		case <-time.After(timeToNext):
			// Proceed with the publishing cycle after a random duration
			f()
		}
	}
}
