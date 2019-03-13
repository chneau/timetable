package timetable

import (
	"errors"
	"sort"
	"time"
)

// Errors ...
var (
	ErrInput      = errors.New("wrong input")
	ErrConstraint = errors.New("constraint fail")
	ErrOverflow   = errors.New("overflow")
)

type point struct {
	time.Time
	val float64
}

// Whener ...
type Whener interface {
	When(time.Time, time.Duration) *time.Time
}

type noopwhener struct{}

func (noopwhener) When(t time.Time, _ time.Duration) *time.Time { return &t }

// TimeTable ...
type TimeTable struct {
	rel        []point
	constraint Whener
	max        float64
}

func check(rel []point, max float64) bool {
	res := 0.
	for _, v := range rel {
		res += v.val
		if res > max {
			return false
		}
	}
	return true
}

func simplify(rel []point) []point {
	new := []point{}
	for i := range rel {
		if len(new) == 0 {
			new = append(new, rel[i])
			continue
		}
		if rel[i].Time.Equal(new[len(new)-1].Time) {
			new[len(new)-1].val += rel[i].val
			if new[len(new)-1].val == 0 {
				new = new[: len(new)-1 : len(new)-1]
			}
		} else {
			new = append(new, rel[i])
		}
	}
	return new
}

func sortPoints(x []point) {
	sort.Slice(x, func(i, j int) bool {
		if x[i].Equal(x[j].Time) {
			return x[i].val < x[j].val
		}
		return x[i].Before(x[j].Time)
	})
}

func (tt *TimeTable) check(from time.Time, dur time.Duration, cap float64) ([]point, bool) {
	a := point{Time: from, val: cap}
	b := point{Time: from.Add(dur), val: -cap}
	x := append(append(tt.rel[:0:0], tt.rel...), a, b)
	sortPoints(x)
	return x, check(x, tt.max)
}

// Add will add the time else returns an error
func (tt *TimeTable) Add(from time.Time, dur time.Duration, cap float64) error {
	if cap < 0 {
		return ErrInput
	}
	if cap > tt.max {
		return ErrInput
	}
	if d := tt.constraint.When(from, dur); d != nil {
		if !from.Equal(*d) {
			return ErrConstraint
		}
	} else {
		return ErrConstraint
	}
	x, ok := tt.check(from, dur, cap)
	if !ok {
		return ErrOverflow
	}
	tt.rel = simplify(x)
	return nil
}

// When returns the soonest time from "from" that satisfies constraints
// else it will return a nil pointer
func (tt *TimeTable) When(from time.Time, dur time.Duration, cap float64) *time.Time {
	// check once
	if t := tt.constraint.When(from, dur); t != nil {
		from = *t
	}
	_, ok := tt.check(from, dur, cap)
	if ok {
		return &from
	}
	for i := range tt.rel {
		if tt.rel[i].After(from) {
			from = tt.rel[i].Time
		} else {
			continue
		}
		if t := tt.constraint.When(from, dur); t != nil {
			from = *t
		}
		_, ok := tt.check(from, dur, cap)
		if ok {
			return &from
		}
	}
	return nil
}

// New ...
func New(max float64, nd Whener) *TimeTable {
	if nd == nil {
		nd = noopwhener{}
	}
	t := &TimeTable{
		max:        max,
		constraint: nd,
	}
	return t
}
