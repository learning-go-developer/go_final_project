package api

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const dateLayout = "20060102"

func afterNow(date, now time.Time) bool {
	y1, m1, d1 := date.Date()
	y2, m2, d2 := now.Date()

	left := time.Date(y1, m1, d1, 0, 0, 0, 0, time.UTC)
	right := time.Date(y2, m2, d2, 0, 0, 0, 0, time.UTC)

	return left.After(right)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeat is empty")
	}

	date, err := time.Parse(dateLayout, dstart)
	if err != nil {
		return "", err
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("invalid repeat")
	}

	switch parts[0] {

	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}

	case "d":
		if len(parts) != 2 {
			return "", errors.New("invalid repeat")
		}

		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", err
		}

		if days < 1 || days > 400 {
			return "", errors.New("invalid interval")
		}

		for {
			date = date.AddDate(0, 0, days)
			if afterNow(date, now) {
				break
			}
		}

	default:
		return "", errors.New("unsupported repeat format")
	}

	return date.Format(dateLayout), nil
}
