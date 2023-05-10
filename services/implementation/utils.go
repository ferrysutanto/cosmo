package implementation

import (
	"math"
	"time"
)

func getBeginningDateOfLastMonth(d time.Time) time.Time {
	// get the beginning of last month with time 00:00:00
	return time.Date(d.Year(), d.Month()-1, 1, 0, 0, 0, 0, d.Location())
}

func getFinalDateOfLastMonth(d time.Time) time.Time {
	// get the final date of last month with time 23:59:59
	return time.Date(d.Year(), d.Month(), 0, 23, 59, 59, 0, d.Location())
}

func getBeginningDateOfTheMonth(d time.Time) time.Time {
	// get the beginning of the month with time 00:00:00
	return time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, d.Location())
}

func getEndOfTheDate(d time.Time) time.Time {
	// get the end of the date with time 23:59:59
	return time.Date(d.Year(), d.Month(), d.Day(), 23, 59, 59, 0, d.Location())
}

func getBeginningOfNextDate(d time.Time) time.Time {
	// get the beginning of next date with time 00:00:00
	return time.Date(d.Year(), d.Month(), d.Day()+1, 0, 0, 0, 0, d.Location())
}

func getEquivalentDateLastMonth(d time.Time) time.Time {
	// get the same date from last month with time 00:00:00, and if last month total date is less than current date, then return the last date of last month
	lastMonthDate := time.Date(d.Year(), d.Month()-1, d.Day(), 0, 0, 0, 0, d.Location())

	if lastMonthDate.Day() < d.Day() {
		return time.Date(d.Year(), d.Month(), 0, 0, 0, 0, 0, d.Location())
	}

	return lastMonthDate
}

func roundTo2DecimalPlaces(f float64) float64 {
	resp := math.Round(f*100) / 100
	return resp
}

func getPercentageChange(old, new float64) float64 {
	if old == 0 {
		return new * 100
	}

	resp := ((new - old) / old) * 100

	if math.IsNaN(resp) {
		return 0
	}

	resp = roundTo2DecimalPlaces(resp)

	return resp
}
