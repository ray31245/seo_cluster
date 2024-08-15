package publishManager

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_randomTime(t *testing.T) {
	tests := []struct {
		name string
		// want int32
	}{
		{name: "test1"},
		{name: "test2"},
		{name: "test3"},
		{name: "test4"},
		{name: "test5"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			count := 10000
			var total time.Duration = 0
			for i := 0; i < count; i++ {
				res := randomTime()
				total += res
			}
			avg := total / time.Duration(count)
			// valid avg time is around 864 minutes
			assert.True(avg > 834*time.Minute && avg < 894*time.Minute)
		})
	}
}

func Test_randomNum(t *testing.T) {
	tests := []struct {
		name string
		want int32
	}{
		{name: "test1"},
		{name: "test2"},
		{name: "test3"},
		{name: "test4"},
		{name: "test5"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			count := 10000
			var total uint64 = 0
			for i := 0; i < count; i++ {
				res := randomNum()
				total += uint64(res)
			}
			var avg float64 = float64(total) * 10 / float64(count)
			// valid avg time is around 5
			assert.True(avg > 4.5 && avg < 5.5)
		})
	}
}
