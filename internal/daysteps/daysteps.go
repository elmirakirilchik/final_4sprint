package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	stepLength = 0.65
	mInKm      = 1000
)

func parsePackage(data string) (int, time.Duration, error) {

	if data == "" {
		return 0, 0, errors.New("invalid")
	}

	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, errors.New("invalid")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil || steps <= 0 {
		return 0, 0, errors.New("invalid")
	}

	durStr := parts[1]

	if strings.ContainsAny(durStr, " -") {
		return 0, 0, errors.New("invalid")
	}

	duration, err := time.ParseDuration(durStr)
	if err != nil || duration <= 0 {
		return 0, 0, errors.New("invalid")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {

	steps, _, err := parsePackage(data)
	if err != nil {
		return ""
	}

	distance := float64(steps) * stepLength / mInKm

	// формула калорий

	calories := distance * weight * 0.605

	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, distance, calories,
	)
}
