package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("invalid format")
	}

	// steps
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, err
	}
	if steps <= 0 {
		return 0, "", 0, errors.New("steps must be greater than 0")
	}

	// training type
	trainingType := parts[1]

	// duration
	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, err
	}

	return steps, trainingType, duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	distanceMeters := float64(steps) * stepLength
	return distanceMeters / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	dist := distance(steps, height)
	return dist / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trainingType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	var calories float64

	switch trainingType {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}

	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}

	default:
		return "", errors.New("unknown training type")
	}

	info := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
		trainingType,
		duration.Hours(),
		dist,
		speed,
		calories,
	)

	return info, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("steps must be greater than 0")
	}
	if weight <= 0 {
		return 0, errors.New("weight must be greater than 0")
	}
	if height <= 0 {
		return 0, errors.New("height must be greater than 0")
	}
	if duration <= 0 {
		return 0, errors.New("duration must be greater than 0")
	}

	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()

	calories := weight * speed * durationInMinutes / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("steps must be greater than 0")
	}
	if weight <= 0 {
		return 0, errors.New("weight must be greater than 0")
	}
	if height <= 0 {
		return 0, errors.New("height must be greater than 0")
	}
	if duration <= 0 {
		return 0, errors.New("duration must be greater than 0")
	}

	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()

	calories := weight * speed * durationInMinutes / minInH * walkingCaloriesCoefficient

	return calories, nil
}
