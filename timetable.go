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

// Point represent a moment in time
type Point struct {
	Time time.Time
	Val  float64
}

// Whener ...
type Whener interface {
	When(time.Time, time.Duration) *time.Time
}

// NoopWhen ...
type NoopWhen struct{}

// When ...
func (NoopWhen) When(t time.Time, _ time.Duration) *time.Time { return &t }

// TimeTable ...
type TimeTable struct {
	Rel        []Point
	Constraint Whener
	Max        float64
}

func check(rel []Point, max float64) bool {
	res := 0.
	for _, v := range rel {
		res += v.Val
		if res > max {
			return false
		}
	}
	return true
}

func simplify(rel []Point) []Point {
	new := []Point{}
	for i := range rel {
		if len(new) == 0 {
			new = append(new, rel[i])
			continue
		}
		if rel[i].Time.Equal(new[len(new)-1].Time) {
			new[len(new)-1].Val += rel[i].Val
			if new[len(new)-1].Val == 0 {
				new = new[: len(new)-1 : len(new)-1]
			}
		} else {
			new = append(new, rel[i])
		}
	}
	return new
}

func sortPoints(x []Point) {
	sort.Slice(x, func(i, j int) bool {
		if x[i].Time.Equal(x[j].Time) {
			return x[i].Val < x[j].Val
		}
		return x[i].Time.Before(x[j].Time)
	})
}

func (tt *TimeTable) check(from time.Time, dur time.Duration, cap float64) ([]Point, bool) {
	a := Point{Time: from, Val: cap}
	b := Point{Time: from.Add(dur), Val: -cap}
	x := append(append(tt.Rel[:0:0], tt.Rel...), a, b)
	sortPoints(x)
	return x, check(x, tt.Max)
}

// Add will add the time else returns an error
func (tt *TimeTable) Add(from time.Time, dur time.Duration, cap float64) error {
	if cap < 0 {
		return ErrInput
	}
	if cap > tt.Max {
		return ErrInput
	}
	if d := tt.Constraint.When(from, dur); d != nil {
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
	tt.Rel = simplify(x)
	return nil
}

// Clone ...
func (tt TimeTable) Clone() TimeTable { return tt }

// Merge ...
func (tt TimeTable) Merge(other TimeTable) *TimeTable {
	tt.Rel = append(tt.Rel, other.Rel...)
	tt.Rel = simplify(tt.Rel)
	if check(tt.Rel, tt.Max) {
		return &tt
	}
	return nil
}

// When returns the soonest time from "from" that satisfies constraints
// else it will return a nil Pointer
func (tt *TimeTable) When(from time.Time, dur time.Duration, cap float64) *time.Time {
	// check once
	if t := tt.Constraint.When(from, dur); t != nil {
		from = *t
	}
	_, ok := tt.check(from, dur, cap)
	if ok {
		return &from
	}
	for i := range tt.Rel {
		if tt.Rel[i].Time.After(from) {
			from = tt.Rel[i].Time
		} else {
			continue
		}
		if t := tt.Constraint.When(from, dur); t != nil {
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
		nd = NoopWhen{}
	}
	t := &TimeTable{
		Max:        max,
		Constraint: nd,
	}
	return t
}
