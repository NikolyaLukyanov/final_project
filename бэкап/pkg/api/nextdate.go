package api

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"
const maxDays = 400

var (
	ErrEmptyRepeat      = errors.New("empty repeat rule")
	ErrInvalidFormat    = errors.New("invalid repeat format")
	ErrInvalidDate      = errors.New("invalid date")
	ErrInvalidDay       = errors.New("invalid day")
	ErrInvalidMonth     = errors.New("invalid month")
	ErrInvalidWeekday   = errors.New("invalid weekday")
	ErrIntervalTooLarge = errors.New("interval exceeds maximum")
	ErrUnsupportedRule  = errors.New("unsupported repeat rule")
)

func NextDate(now time.Time, dateStr string, repeat string) (string, error) {
	if repeat == "" {
		return "", ErrEmptyRepeat
	}

	date, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		return "", ErrInvalidDate
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", ErrInvalidFormat
	}

	switch parts[0] {
	case "d":
		return handleDailyRepeat(now, date, parts)
	case "y":
		return handleYearlyRepeat(now, date)
	case "w":
		return handleWeeklyRepeat(now, date, parts)
	case "m":
		return handleMonthlyRepeat(now, date, parts)
	default:
		return "", ErrUnsupportedRule
	}
}

func handleDailyRepeat(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", ErrInvalidFormat
	}

	interval, err := strconv.Atoi(parts[1])
	if err != nil || interval <= 0 || interval > maxDays {
		return "", ErrIntervalTooLarge
	}

	for {
		date = date.AddDate(0, 0, interval)
		if afterNow(date, now) {
			break
		}
	}

	return date.Format(dateFormat), nil
}

func handleYearlyRepeat(now, date time.Time) (string, error) {
	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			break
		}
	}

	return date.Format(dateFormat), nil
}

func handleWeeklyRepeat(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", ErrInvalidFormat
	}

	daysStr := strings.Split(parts[1], ",")
	weekdays := make(map[int]bool)
	for _, dayStr := range daysStr {
		day, err := strconv.Atoi(dayStr)
		if err != nil || day < 1 || day > 7 {
			return "", ErrInvalidWeekday
		}
		weekdays[day] = true
	}

	for {
		date = date.AddDate(0, 0, 1)
		if afterNow(date, now) {
			weekday := int(date.Weekday())
			if weekday == 0 {
				weekday = 7 // Sunday
			}
			if weekdays[weekday] {
				break
			}
		}
	}

	return date.Format(dateFormat), nil
}

func handleMonthlyRepeat(now, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", ErrInvalidFormat
	}

	daysStr := strings.Split(parts[1], ",")
	days := make(map[int]bool)
	for _, dayStr := range daysStr {
		day, err := strconv.Atoi(dayStr)
		if err != nil {
			return "", ErrInvalidDay
		}
		if day == -1 || day == -2 {
			days[day] = true
			continue
		}
		if day < 1 || day > 31 {
			return "", ErrInvalidDay
		}
		days[day] = true
	}

	var months map[int]bool
	if len(parts) > 2 {
		monthsStr := strings.Split(parts[2], ",")
		months = make(map[int]bool)
		for _, monthStr := range monthsStr {
			month, err := strconv.Atoi(monthStr)
			if err != nil || month < 1 || month > 12 {
				return "", ErrInvalidMonth
			}
			months[month] = true
		}
	}

	for {
		date = date.AddDate(0, 0, 1)
		if afterNow(date, now) {
			day := date.Day()
			month := int(date.Month())

			// Check for last day or penultimate day
			lastDay := daysInMonth(date.Year(), month)
			if days[-1] && day == lastDay {
				break
			}
			if days[-2] && day == lastDay-1 {
				break
			}

			// Check regular days
			if (months == nil || months[month]) && days[day] {
				break
			}
		}
	}

	return date.Format(dateFormat), nil
}

func daysInMonth(year, month int) int {
	return time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func afterNow(date, now time.Time) bool {
	return date.After(now)
}
