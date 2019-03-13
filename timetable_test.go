package timetable

import (
	"testing"
	"time"

	"github.com/chneau/openhours"
)

func TestTimeTable_Add(t *testing.T) {
	t.Run("overlapping and simplifying", func(t *testing.T) {
		l := time.Now().Location()
		oh := openhours.New("mo-fr 11:00-16:00", l)
		tt := New(5, oh)
		for i := 0; i < 5; i++ {
			err := tt.Add(time.Date(2019, 3, 12, 11, 0, 0, 0, l), time.Hour*2, 1)
			if err != nil {
				t.Error(err)
			}
			err = tt.Add(time.Date(2019, 3, 12, 13, 0, 0, 0, l), time.Hour*2, 1)
			if err != nil {
				t.Error(err)
			}
		}
		err := tt.Add(time.Date(2019, 3, 12, 11, 0, 0, 0, l), time.Hour*2, 1)
		if err == nil {
			t.Error("no error is no good at this point")
		}
	})
	t.Run("ranges overlap at same time", func(t *testing.T) {
		l := time.Now().Location()
		oh := openhours.New("mo-fr 11:00-16:00", l)
		tt := New(2, oh)
		err := tt.Add(time.Date(2019, 3, 12, 11, 0, 0, 0, l), time.Hour, 1)
		if err != nil {
			t.Error(err)
		}
		err = tt.Add(time.Date(2019, 3, 12, 12, 0, 0, 0, l), time.Hour, 2)
		if err != nil {
			t.Error(err)
		}
		err = tt.Add(time.Date(2019, 3, 12, 13, 0, 0, 0, l), time.Hour, 1)
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("micro overflow", func(t *testing.T) {
		l := time.Now().Location()
		oh := openhours.New("mo-fr 11:00-16:00", l)
		tt := New(1, oh)
		err := tt.Add(time.Date(2019, 3, 12, 11, 0, 0, 0, l), time.Hour*5, 1)
		if err != nil {
			t.Error(err)
		}
		err = tt.Add(time.Date(2019, 3, 12, 12, 0, 0, 0, l), time.Microsecond, 0.0001)
		if err == nil {
			t.Error("no error is no good at this point")
		}
	})
}

func TestTimeTable_When(t *testing.T) {
	t.Error("Must continue implementation")
}
