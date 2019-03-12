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

// Add ...
func (tt *TimeTable) Add(from time.Time, dur time.Duration, cap float64) error {
	if cap < 0 {
		return ErrInput
	}
	if cap > tt.max {
		return ErrInput
	}
	if tt.constraint != nil {
		if d := tt.constraint.When(from, dur); d != nil {
			if !from.Equal(*d) {
				return ErrConstraint
			}
		} else {
			return ErrConstraint
		}
	}
	a := point{Time: from, val: cap}
	b := point{Time: from.Add(dur), val: -cap}
	x := append(append(tt.rel[:0:0], tt.rel...), a, b)
	sortPoints(x)
	if !check(x, tt.max) {
		return ErrOverflow
	}
	tt.rel = simplify(x)
	return nil
}

// When ...
func (tt *TimeTable) When(t time.Time, d time.Duration, cap float64) *time.Time {
	score := 0.
	if tt.constraint != nil {
		if newT := tt.constraint.When(t, d); newT != nil {
			if !newT.Equal(t) {
				t = *newT
			}
		}
	}
	for i := range tt.rel {
		score += tt.rel[i].val
		if tt.rel[i].Before(t) {
			continue
		}
		if i+1 < len(tt.rel) {

		}
	}
	return nil
}

// New ...
func New(max float64, nd Whener) *TimeTable {
	t := &TimeTable{
		max:        max,
		constraint: nd,
	}
	return t
}
