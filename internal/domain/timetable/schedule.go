package timetable

import (
	"sort"
	"time"
)

// WeekdaysToMask packs ISO weekdays (1..7) into a bitmask, bit (weekday-1).
func WeekdaysToMask(days []int) int {
	mask := 0
	for _, d := range days {
		if d >= 1 && d <= 7 {
			mask |= 1 << (d - 1)
		}
	}
	return mask
}

// MaskToWeekdays unpacks a bitmask back into sorted ISO weekdays.
func MaskToWeekdays(mask int) []int {
	out := make([]int, 0, 7)
	for d := 1; d <= 7; d++ {
		if mask&(1<<(d-1)) != 0 {
			out = append(out, d)
		}
	}
	return out
}

// isoWeekday converts Go's time.Weekday (Sun=0 … Sat=6) to ISO (Mon=1 … Sun=7).
func isoWeekday(t time.Time) int {
	if wd := int(t.Weekday()); wd != 0 {
		return wd
	}
	return 7
}

func truncateDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// GenerateSessions expands schedule rules into concrete class occurrences that
// fall within [from, to], each clamped to its batch's start/end dates.
func GenerateSessions(rules []Timetable, from, to time.Time) []ClassSession {
	from = truncateDay(from)
	to = truncateDay(to)

	var out []ClassSession
	for _, r := range rules {
		start := from
		if bs := truncateDay(r.BatchStart); bs.After(start) {
			start = bs
		}
		end := to
		if be := truncateDay(r.BatchEnd); be.Before(end) {
			end = be
		}

		wanted := make(map[int]bool, len(r.Weekdays))
		for _, d := range r.Weekdays {
			wanted[d] = true
		}

		for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
			if !wanted[isoWeekday(d)] {
				continue
			}
			out = append(out, ClassSession{
				Date:        d,
				Weekday:     isoWeekday(d),
				StartTime:   r.StartTime,
				EndTime:     r.EndTime,
				CourseID:    r.CourseID,
				CourseTitle: r.CourseTitle,
				BatchID:     r.BatchID,
				BatchName:   r.BatchName,
			})
		}
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].Date.Equal(out[j].Date) {
			return out[i].StartTime < out[j].StartTime
		}
		return out[i].Date.Before(out[j].Date)
	})

	return out
}
