package timetable

import (
	"testing"
	"time"

	"github.com/chneau/openhours"
)

func TestTimeTable_Add(t *testing.T) {
	t.Run("overlapping and simplifying", func(t *testing.T) {
		l := time.Local
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
		l := time.Local
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
		l := time.Local
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
		l := time.Local
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
		l := time.Local
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

func BenchmarkTimeTable(b *testing.B) {
	oh := openhours.New("Mo-Fr 08:00-18:00", time.UTC)
	baseTime := time.Date(2026, 8, 24, 9, 0, 0, 0, time.UTC)
	oneHour := time.Hour

	b.Run("Sequential Add", func(b *testing.B) {
		tt := New(100.0, oh)
		for i := 0; i < b.N; i++ {
			_ = tt.Add(baseTime.Add(time.Duration(i%500)*10*time.Minute), oneHour, 1.0)
			if len(tt.Rel) > 200 {
				tt = New(100.0, oh)
			}
		}
	})

	b.Run("Capacity Search When", func(b *testing.B) {
		tt := New(3.0, oh)
		for i := 0; i < 5; i++ {
			_ = tt.Add(baseTime.Add(time.Duration(i*2)*time.Hour), 2*time.Hour, 2.0)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			tt.When(baseTime.Add(time.Duration(i%40)*time.Hour), oneHour, 2.0)
		}
	})
}

func TestRunStandardBenchmarkSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}
	oh := openhours.New("Mo-Fr 08:00-18:00", time.UTC)
	baseTime := time.Date(2026, 8, 24, 9, 0, 0, 0, time.UTC)
	oneHour := time.Hour

	println("========================================================")
	println("Running TimeTable Benchmarks (Go Suite)")
	println("========================================================")

	// Warmup
	warmupTT := New(5.0, oh)
	for i := 0; i < 5000; i++ {
		_ = warmupTT.Add(baseTime.Add(time.Duration(i%24)*time.Hour), oneHour, 1.0)
		warmupTT.When(baseTime.Add(time.Duration(i%10)*time.Hour), oneHour, 1.0)
		if len(warmupTT.Rel) > 50 {
			warmupTT = New(5.0, oh)
		}
	}

	// 1. Sequential Add
	addOps := 50000
	tt := New(100.0, oh)
	t0 := time.Now()
	for i := 0; i < addOps; i++ {
		_ = tt.Add(baseTime.Add(time.Duration(i*10)*time.Minute), oneHour, 1.0)
		if len(tt.Rel) > 200 {
			tt = New(100.0, oh)
		}
	}
	d1 := time.Since(t0)
	println("1. Sequential Add (50000 calls):              ", d1.Milliseconds(), "ms", "(", float64(d1.Microseconds())/float64(addOps), "us/op )")

	// 2. Capacity Search When
	whenOps := 20000
	searchTT := New(3.0, oh)
	for i := 0; i < 5; i++ {
		_ = searchTT.Add(baseTime.Add(time.Duration(i*2)*time.Hour), 2*time.Hour, 2.0)
	}
	t0 = time.Now()
	for i := 0; i < whenOps; i++ {
		searchTT.When(baseTime.Add(time.Duration(i%40)*time.Hour), oneHour, 2.0)
	}
	d2 := time.Since(t0)
	println("2. Capacity Search When (20000 calls):         ", d2.Milliseconds(), "ms", "(", float64(d2.Microseconds())/float64(whenOps), "us/op )")
	println("========================================================")
}
