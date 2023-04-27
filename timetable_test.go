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
		tt := New(1., oh)
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
	t.Run("micro overflow", func(t *testing.T) {
		l := time.Now().Location()
		oh := openhours.New("mo-fr 11:00-16:00", l) // this is the slow part of the code there this one
		tt := New(10, oh)
		d := time.Date(2019, 3, 12, 10, 0, 0, 0, l)
		for i := 0; i < 1000; i++ {
			when := tt.When(d, time.Hour, 1)
			if when.After(d) { // when must never be null
				d = *when // makes test faster and more realistic ?
			}
			err := tt.Add(*when, time.Hour, 1)
			if err != nil {
				t.Error("this should not be failing")
			}
		}
	})
	t.Run("should be nil", func(t *testing.T) {
		l := time.Now().Location()
		oh := openhours.New("mo-fr 11:00-16:00", l)
		tt := New(10, oh)
		d := time.Date(2019, 3, 12, 10, 0, 0, 0, l)
		when := tt.When(d, time.Hour, 12)
		if when != nil {
			t.Error("when should be nil")
		}
	})
}

func TestTimeTable_Merge(t *testing.T) {
	t.Run("merge success", func(t *testing.T) {
		l := time.Local
		oh := openhours.New("mo-fr 11:00-16:00", l)
		tt := New(10, oh)
		d := time.Date(2019, 3, 12, 10, 0, 0, 0, l)
		when := tt.When(d, time.Hour, 4)
		err := tt.Add(*when, time.Hour, 4)
		if err != nil {
			t.Error(err)
			return
		}
		tt2 := tt.Clone()
		tt3 := tt.Merge(tt2)
		if tt3 == nil {
			t.Error("merge failed")
		}
	})
	t.Run("merge fail", func(t *testing.T) {
		l := time.Local
		oh := openhours.New("mo-fr 11:00-16:00", l)
		tt := New(10, oh)
		d := time.Date(2019, 3, 12, 10, 0, 0, 0, l)
		when := tt.When(d, time.Hour, 6)
		err := tt.Add(*when, time.Hour, 6)
		if err != nil {
			t.Error(err)
			return
		}
		tt2 := tt.Clone()
		tt3 := tt.Merge(tt2)
		if tt3 != nil {
			t.Error("merge should have failed")
		}
	})
}
