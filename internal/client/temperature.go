package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

type AirConditionerMode int

const (
	ModeAuto AirConditionerMode = iota + 1
	ModeCool
	ModeDry
	ModeFan
	ModeHeat
)

func (c Client) GetNowTemperature(ctx context.Context, deviceID string) (float64, error) {
	path := fmt.Sprintf("/devices/%s/status", deviceID)
	res, err := c.get(ctx, path)
	if err != nil {
		return 0, fmt.Errorf("failed get device status %w", err)
	}
	var data struct {
		StatusCode int `json:"statusCode"`
		Body       struct {
			Temperature float64 `json:"temperature"`
		} `json:"body"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(res, &data); err != nil {
		return 0, err
	}
	if data.StatusCode != 200 {
		return 0, fmt.Errorf("unexpected status code: %d, message: %s", data.StatusCode, data.Message)
	}
	return data.Body.Temperature, nil
}

func (c Client) SetTemperature(ctx context.Context, deviceID string, temperature float64, mode AirConditionerMode) error {
	path := fmt.Sprintf("/devices/%s/commands", deviceID)
	payload := struct {
		Command     string `json:"command"`
		Parameter   string `json:"parameter"`
		CommandType string `json:"commandType"`
	}{
		Command:     "setTemperature",
		Parameter:   fmt.Sprintf("%.1f,%d,1,on", temperature, mode),
		CommandType: "command",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed marshal payload: %w", err)
	}
	_, err = c.post(ctx, path, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed post set temperature command: %w", err)
	}
	return nil
}
