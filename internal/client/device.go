package client

import (
	"context"
	"encoding/json"
	"fmt"
)

const (
	MeterPro       = "MeterPro"
	AirConditioner = "Air Conditioner"
)

type ListDeviceResponse struct {
	StatusCode int            `json:"statusCode"`
	Body       listDeviceBody `json:"body"`
	Message    string         `json:"message"`
}

type listDeviceBody struct {
	DeviceList         []device         `json:"deviceList"`
	InfraredRemoteList []infraredRemote `json:"infraredRemoteList"`
}

type device struct {
	DeviceID           string `json:"deviceId"`
	DeviceName         string `json:"deviceName"`
	DeviceType         string `json:"deviceType"`
	EnableCloudService bool   `json:"enableCloudService"`
	HubDeviceID        string `json:"hubDeviceId"`
}

type infraredRemote struct {
	DeviceID    string `json:"deviceId"`
	DeviceName  string `json:"deviceName"`
	RemoteType  string `json:"remoteType"`
	HubDeviceID string `json:"hubDeviceId"`
}

func (c *Client) listDevice(ctx context.Context) (*ListDeviceResponse, error) {
	path := "/devices"
	res, err := c.get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed get device list %w", err)
	}
	var data ListDeviceResponse
	if err := json.Unmarshal(res, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (c Client) MultiGetAirConditioners(ctx context.Context) ([]*infraredRemote, error) {
	var res []*infraredRemote
	devices, err := c.listDevice(ctx)
	if err != nil {
		return nil, err
	}
	for _, device := range devices.Body.InfraredRemoteList {
		if device.RemoteType == AirConditioner {
			res = append(res, &device)
		}
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("air conditioner not found")
	}
	return res, nil
}

func (c Client) GetMeterPro(ctx context.Context) (*device, error) {
	devices, err := c.listDevice(ctx)
	if err != nil {
		return nil, err
	}
	for _, device := range devices.Body.DeviceList {
		if device.DeviceType == MeterPro {
			return &device, nil
		}
	}
	return nil, fmt.Errorf("meter pro not found")
}
