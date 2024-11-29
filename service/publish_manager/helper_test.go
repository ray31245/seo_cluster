package publishmanager

import (
	"context"
	"log"
	"reflect"
	"testing"
	"time"
)

func Test_computeTimePointArray(t *testing.T) {
	tests := []struct {
		name string
		want []time.Time
	}{
		// TODO: Add test cases.
		{name: "Test 1", want: []time.Time{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := computeTimePointArray()
			log.Println(len(got))
			if got := computeTimePointArray(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("computeTimePointArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_timeArrSchedule(t *testing.T) {
	type args struct {
		ctx     context.Context
		timeArr []time.Time
		f       func()
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "Test 1", args: args{
			ctx: context.Background(),
			timeArr: []time.Time{
				time.Now().Add(1 * time.Second),
				time.Now().Add(3 * time.Second),
				time.Now().Add(15 * time.Second),
			},
			f: func() { log.Println("Hello") },
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Println(tt.args.timeArr)
			timeArrSchedule(tt.args.ctx, tt.args.timeArr, tt.args.f)
		})
	}
}
