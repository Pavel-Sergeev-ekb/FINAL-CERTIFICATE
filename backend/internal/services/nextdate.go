package services

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
	ErrUnsupported    = errors.New("unsupported repeat rule")
)

func NextDate(now time.Time, dStart string, repeat string) (string, error) {

	// Пустой repeat
	if repeat == "" {
		return "", ErrBadRepeat
	}

	startDate, err := time.Parse(TypeDay, dStart)
	if err != nil {
		return "", ErrInvalidDate
	}

	parts := strings.Split(repeat, " ")
	rule := parts[0]

	switch rule {
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid daily rule format")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 || days > 400 {
			return "", fmt.Errorf("invalid number of days")
		}

		for {
			startDate = startDate.AddDate(0, 0, days)
			if startDate.After(now) {
				break
			}
		}

	case "y":
		for {
			startDate = startDate.AddDate(1, 0, 0)
			if startDate.After(now) {
				break
			}
		}

	default:
		return "", ErrUnsupported
	}

	return startDate.Format(TypeDay), nil
}
