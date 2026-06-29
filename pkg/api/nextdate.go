package api

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const dateLayout = "20060102"

func NextDate(now time.Time, start string, repeat string) (string, error) {
	cleanRule := strings.TrimSpace(repeat)
	if cleanRule == "" {
		return "", errors.New("empty repeat rule")
	}

	startDate, err := time.Parse(dateLayout, start)
	if err != nil {
		return "", err
	}

	current := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
	refTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	chunks := strings.Fields(cleanRule)
	if len(chunks) == 0 {
		return "", errors.New("invalid repeat format")
	}

	switch chunks[0] {
	case "y":
		for {
			current = current.AddDate(1, 0, 0)
			if current.After(refTime) {
				break
			}
		}

	case "d":
		if len(chunks) != 2 {
			return "", errors.New("missing days interval")
		}
		step, err := strconv.Atoi(chunks[1])
		if err != nil || step < 1 || step > 400 {
			return "", errors.New("incorrect days interval")
		}

		for {
			current = current.AddDate(0, 0, step)
			if current.After(refTime) {
				break
			}
		}

	default:
		return "", errors.New("unknown rule type")
	}

	return current.Format(dateLayout), nil
}
