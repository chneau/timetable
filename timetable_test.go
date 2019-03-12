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
}

func TestTimeTable_When(t *testing.T) {
	t.Error("Must continue implementation")
}
