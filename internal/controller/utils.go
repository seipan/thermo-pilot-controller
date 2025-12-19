package controller

import (
	"fmt"
	"strconv"
)

func ParseTemperature(temp string) (float64, error) {
	value, err := strconv.ParseFloat(temp, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid temperature format: %s", temp)
	}
	return value, nil
}

func FormatTemperature(temp float64) string {
	return fmt.Sprintf("%.1f", temp)
}
