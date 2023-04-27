package timetable

import (
	"errors"
	"sort"
	"time"

	"golang.org/x/exp/constraints"
)

// Errors ...
var (
	ErrInput      = errors.New("wrong input")
	ErrConstraint = errors.New("constraint fail")
	ErrOverflow   = errors.New("overflow")
)

type Number interface {
	constraints.Float | constraints.Integer
}

// Point represent a moment in time
type Point[T Number] struct {
	Time time.Time
	Val  T
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
type TimeTable[T Number] struct {
	Rel        []Point[T]
	Constraint Whener
	Max        T
}

func check[T Number](rel []Point[T], max T) bool {
	var res T
	for _, v := range rel {
		res += v.Val
		if res > max {
			return false
		}
	}
	return true
}

func simplify[T Number](rel []Point[T]) []Point[T] {
	new := []Point[T]{}
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

func sortPoints[T Number](x []Point[T]) {
	sort.Slice(x, func(i, j int) bool {
		if x[i].Time.Equal(x[j].Time) {
			return x[i].Val < x[j].Val
		}
		return x[i].Time.Before(x[j].Time)
	})
}

func (tt *TimeTable[T]) check(from time.Time, dur time.Duration, cap T) ([]Point[T], bool) {
	a := Point[T]{Time: from, Val: cap}
	b := Point[T]{Time: from.Add(dur), Val: -cap}
	x := append(append(tt.Rel[:0:0], tt.Rel...), a, b)
	sortPoints(x)
	return x, check(x, tt.Max)
}

// Add will add the time else returns an error
func (tt *TimeTable[T]) Add(from time.Time, dur time.Duration, cap T) error {
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
func (tt TimeTable[T]) Clone() TimeTable[T] { return tt }

// Merge ...
func (tt TimeTable[T]) Merge(other TimeTable[T]) *TimeTable[T] {
	tt.Rel = append(tt.Rel, other.Rel...)
	tt.Rel = simplify(tt.Rel)
	if check(tt.Rel, tt.Max) {
		return &tt
	}
	return nil
}

// When returns the soonest time from "from" that satisfies constraints
// else it will return a nil Pointer
func (tt *TimeTable[T]) When(from time.Time, dur time.Duration, cap T) *time.Time {
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
func New[T Number](max T, nd Whener) *TimeTable[T] {
	if nd == nil {
		nd = NoopWhen{}
	}
	t := &TimeTable[T]{
		Max:        max,
		Constraint: nd,
	}
	return t
}
