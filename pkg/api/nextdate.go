package api

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const dateLayout = "20060102"

func NextDate(now time.Time, start string, repeat string) (string, error) {
	if strings.TrimSpace(repeat) == "" {
		return "", errors.New("empty repeat rule")
	}

	current, err := time.Parse(dateLayout, start)
	if err != nil {
		return "", err
	}

	current = time.Date(current.Year(), current.Month(), current.Day(), 0, 0, 0, 0, time.UTC)
	refTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	args := strings.Fields(repeat)
	if len(args) == 0 {
		return "", errors.New("invalid repeat format")
	}

	switch args[0] {
	case "y":
		for {
			current = current.AddDate(1, 0, 0)
			if current.After(refTime) {
				break
			}
		}

	case "d":
		if len(args) != 2 {
			return "", errors.New("missing days interval")
		}
		days, err := strconv.Atoi(args[1])
		if err != nil || days < 1 || days > 400 {
			return "", errors.New("incorrect days interval")
		}

		current = current.AddDate(0, 0, days)

		for !current.After(refTime) {
			current = current.AddDate(0, 0, days)
		}

	default:
		return "", errors.New("unknown rule type")
	}

	return current.Format(dateLayout), nil
}
