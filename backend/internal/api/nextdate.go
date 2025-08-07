package api

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const TypeDay = "20060102"

var (
	ErrBadRepeat      = errors.New("invalid missing repeat rule")
	ErrInvalidDate    = errors.New("invalid date format, expected YYYYMMDD")
	ErrIntervalTooBig = errors.New("day interval exceeds max of 400")
	ErrUnsupported    = errors.New("unsupported repeat rule (only d <N> and y)")
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// Парсим dstart
	start, err := time.Parse(TypeDay, dstart)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidDate, err)
	}

	// Пустой repeat
	if repeat == "" {
		return "", ErrBadRepeat
	}

	var stepFunc func(time.Time) (time.Time, error)

	switch {
	case repeat == "y":
		stepFunc = func(t time.Time) (time.Time, error) {
			return t.AddDate(1, 0, 0), nil
		}
	case len(repeat) > 2 && repeat[:2] == "d ":
		var days int
		_, err := fmt.Sscanf(repeat, "d %d", &days)
		if err != nil {
			return "", ErrBadRepeat
		}
		if days < 1 || days > 400 {
			return "", ErrIntervalTooBig
		}
		stepFunc = func(t time.Time) (time.Time, error) {
			return t.AddDate(0, 0, days), nil
		}
	default:
		return "", ErrUnsupported
	}

	// Цикл для поиска следующей даты
	candidate := start
	for {
		candidate, err = stepFunc(candidate)
		if err != nil {
			return "", err
		}
		if candidate.After(now) {
			break
		}
	}
	return candidate.Format(TypeDay), nil
}

func ParseWeekDays(repeat string) ([]int, error) {
	parts := strings.Split(repeat[2:], ",")

	var result []int
	for _, part := range parts {
		dayStr := strings.TrimSpace(part)
		day, err := strconv.Atoi(dayStr)
		if err != nil || day < 1 || day > 7 {
			return nil, ErrBadRepeat
		}
		result = append(result, day)
	}
	return result, nil
}

func FindNextWeekDay(from time.Time, days []int) (time.Time, error) {
	for _, day := range days {
		if day < 1 || day > 7 {
			return time.Time{}, fmt.Errorf("invalid week day: %d", day)
		}
	}
	daySet := make(map[int]struct{})
	for _, day := range days {
		daySet[day] = struct{}{}
	}
	for {
		from = from.AddDate(0, 0, 1)
		weekDay := int(from.Weekday()) + 1

		if _, exists := daySet[weekDay]; exists {
			return from, nil
		}
	}
}

func ParseMonthDays(repeat string) ([]int, []int, error) {
	parts := strings.Split(repeat[2:], " ")
	var daysOfMonth []int
	var months []int

	dayPart := parts[0]
	dayValues := strings.Split(dayPart, ",")
	for _, val := range dayValues {
		val = strings.TrimSpace(val)
		if val == "-1" {
			daysOfMonth = append(daysOfMonth, -1)
		} else if val == "-2" {
			daysOfMonth = append(daysOfMonth, -2)
		} else {
			day, err := strconv.Atoi(val)
			if err != nil || day < 1 || day > 31 {
				return nil, nil, ErrBadRepeat
			}
			daysOfMonth = append(daysOfMonth, day)
		}
	}

	if len(parts) > 1 {
		monthPart := parts[1]
		monthValues := strings.Split(monthPart, ",")
		for _, val := range monthValues {
			val = strings.TrimSpace(val)
			month, err := strconv.Atoi(val)
			if err != nil || month < 1 || month > 12 {
				return nil, nil, ErrBadRepeat
			}
			months = append(months, month)
		}
	}
	if len(daysOfMonth) == 0 {
		return nil, nil, ErrBadRepeat
	}
	return daysOfMonth, months, nil
}

func findNextMonthDay(from time.Time, daysOfMonth []int, months []int) time.Time {
	for {
		currentMonth := from.Month()
		validMonth := len(months) == 0 || contains(months, int(currentMonth))

		currentDay := from.Day()
		lastDay := time.Date(from.Year(), from.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()

		for _, day := range daysOfMonth {
			var targetDay int
			if day == -1 {
				targetDay = lastDay
			} else if day == -2 {
				targetDay = lastDay - 1
			} else {
				targetDay = day
			}
			if validMonth && currentDay == targetDay {
				return from
			}
		}
		from = from.AddDate(0, 0, 1)
	}
}

func contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
