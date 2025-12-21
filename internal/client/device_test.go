package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetAirConditioner(t *testing.T) {
	tests := []struct {
		name           string
		response       ListDeviceResponse
		statusCode     int
		wantDeviceID   string
		wantDeviceName string
		wantRemoteType string
		wantErr        bool
		wantErrMsg     string
	}{
		{
			name: "success - air conditioner found",
			response: ListDeviceResponse{
				StatusCode: 200,
				Message:    "success",
				Body: listDeviceBody{
					InfraredRemoteList: []infraredRemote{
						{
							DeviceID:    "device1",
							DeviceName:  "Living Room TV",
							RemoteType:  "TV",
							HubDeviceID: "hub1",
						},
						{
							DeviceID:    "device2",
							DeviceName:  "Bedroom AC",
							RemoteType:  AirConditioner,
							HubDeviceID: "hub1",
						},
					},
				},
			},
			statusCode:     200,
			wantDeviceID:   "device2",
			wantDeviceName: "Bedroom AC",
			wantRemoteType: AirConditioner,
			wantErr:        false,
		},
		{
			name: "error - air conditioner not found",
			response: ListDeviceResponse{
				StatusCode: 200,
				Message:    "success",
				Body: listDeviceBody{
					InfraredRemoteList: []infraredRemote{
						{
							DeviceID:    "device1",
							DeviceName:  "Living Room TV",
							RemoteType:  "TV",
							HubDeviceID: "hub1",
						},
					},
				},
			},
			statusCode: 200,
			wantErr:    true,
			wantErrMsg: "air conditioner not found",
		},
		{
			name: "error - empty response",
			response: ListDeviceResponse{
				StatusCode: 200,
				Message:    "success",
				Body: listDeviceBody{
					InfraredRemoteList: []infraredRemote{},
				},
			},
			statusCode: 200,
			wantErr:    true,
			wantErrMsg: "air conditioner not found",
		},
		{
			name:       "error - API error",
			response:   ListDeviceResponse{},
			statusCode: 500,
			wantErr:    true,
			wantErrMsg: "unexpected status code: 500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/v1.1/devices", r.URL.Path)
				assert.Equal(t, http.MethodGet, r.Method)
				assert.NotEmpty(t, r.Header.Get("Authorization"))
				assert.NotEmpty(t, r.Header.Get("sign"))
				assert.NotEmpty(t, r.Header.Get("nonce"))
				assert.NotEmpty(t, r.Header.Get("t"))

				w.WriteHeader(tt.statusCode)
				if tt.statusCode == 200 {
					if err := json.NewEncoder(w).Encode(tt.response); err != nil {
						t.Errorf("Failed to encode response: %v", err)
					}
				}
			}))
			defer server.Close()

			client := NewClient("test-token", "test-secret")
			client.HttpClient = server.Client()
			oldAPI := switchBotAPI
			switchBotAPI = server.URL + "/v1.1"
			defer func() { switchBotAPI = oldAPI }()

			ctx := context.Background()
			got, err := client.MultiGetAirConditioners(ctx)

			if tt.wantErr {
				require.Error(t, err)
				if tt.wantErrMsg != "" {
					assert.Contains(t, err.Error(), tt.wantErrMsg)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				assert.Equal(t, tt.wantDeviceID, got[0].DeviceID)
				assert.Equal(t, tt.wantDeviceName, got[0].DeviceName)
				assert.Equal(t, tt.wantRemoteType, got[0].RemoteType)
			}
		})
	}
}

func TestClient_GetMeterPro(t *testing.T) {
	tests := []struct {
		name           string
		response       ListDeviceResponse
		statusCode     int
		wantDeviceID   string
		wantDeviceName string
		wantDeviceType string
		wantErr        bool
		wantErrMsg     string
	}{
		{
			name: "success - meter pro found",
			response: ListDeviceResponse{
				StatusCode: 200,
				Message:    "success",
				Body: listDeviceBody{
					DeviceList: []device{
						{
							DeviceID:           "device1",
							DeviceName:         "Living Room Thermometer",
							DeviceType:         "Meter",
							EnableCloudService: true,
							HubDeviceID:        "hub1",
						},
						{
							DeviceID:           "device2",
							DeviceName:         "Bedroom Meter Pro",
							DeviceType:         MeterPro,
							EnableCloudService: true,
							HubDeviceID:        "hub1",
						},
					},
				},
			},
			statusCode:     200,
			wantDeviceID:   "device2",
			wantDeviceName: "Bedroom Meter Pro",
			wantDeviceType: MeterPro,
			wantErr:        false,
		},
		{
			name: "error - meter pro not found",
			response: ListDeviceResponse{
				StatusCode: 200,
				Message:    "success",
				Body: listDeviceBody{
					DeviceList: []device{
						{
							DeviceID:           "device1",
							DeviceName:         "Living Room Thermometer",
							DeviceType:         "Meter",
							EnableCloudService: true,
							HubDeviceID:        "hub1",
						},
					},
				},
			},
			statusCode: 200,
			wantErr:    true,
			wantErrMsg: "meter pro not found",
		},
		{
			name: "error - empty device list",
			response: ListDeviceResponse{
				StatusCode: 200,
				Message:    "success",
				Body: listDeviceBody{
					DeviceList: []device{},
				},
			},
			statusCode: 200,
			wantErr:    true,
			wantErrMsg: "meter pro not found",
		},
		{
			name:       "error - API error",
			response:   ListDeviceResponse{},
			statusCode: 401,
			wantErr:    true,
			wantErrMsg: "unexpected status code: 401",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/v1.1/devices", r.URL.Path)
				assert.Equal(t, http.MethodGet, r.Method)

				w.WriteHeader(tt.statusCode)
				if tt.statusCode == 200 {
					if err := json.NewEncoder(w).Encode(tt.response); err != nil {
						t.Errorf("Failed to encode response: %v", err)
					}
				}
			}))
			defer server.Close()

			client := NewClient("test-token", "test-secret")
			client.HttpClient = server.Client()

			oldAPI := switchBotAPI
			switchBotAPI = server.URL + "/v1.1"
			defer func() { switchBotAPI = oldAPI }()

			ctx := context.Background()
			got, err := client.GetMeterPro(ctx)

			if tt.wantErr {
				require.Error(t, err)
				if tt.wantErrMsg != "" {
					assert.Contains(t, err.Error(), tt.wantErrMsg)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				assert.Equal(t, tt.wantDeviceID, got.DeviceID)
				assert.Equal(t, tt.wantDeviceName, got.DeviceName)
				assert.Equal(t, tt.wantDeviceType, got.DeviceType)
			}
		})
	}
}

func TestClient_listDevice(t *testing.T) {
	tests := []struct {
		name       string
		response   ListDeviceResponse
		statusCode int
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "success - mixed devices",
			response: ListDeviceResponse{
				StatusCode: 200,
				Message:    "success",
				Body: listDeviceBody{
					DeviceList: []device{
						{
							DeviceID:   "device1",
							DeviceName: "Test Device",
							DeviceType: "Meter",
						},
					},
					InfraredRemoteList: []infraredRemote{
						{
							DeviceID:   "remote1",
							DeviceName: "Test Remote",
							RemoteType: "TV",
						},
					},
				},
			},
			statusCode: 200,
			wantErr:    false,
		},
		{
			name:       "error - unauthorized",
			response:   ListDeviceResponse{},
			statusCode: 401,
			wantErr:    true,
			wantErrMsg: "unexpected status code: 401",
		},
		{
			name:       "error - internal server error",
			response:   ListDeviceResponse{},
			statusCode: 500,
			wantErr:    true,
			wantErrMsg: "unexpected status code: 500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				if tt.statusCode == 200 {
					if err := json.NewEncoder(w).Encode(tt.response); err != nil {
						t.Errorf("Failed to encode response: %v", err)
					}
				}
			}))
			defer server.Close()

			client := NewClient("test-token", "test-secret")
			client.HttpClient = server.Client()

			oldAPI := switchBotAPI
			switchBotAPI = server.URL + "/v1.1"
			defer func() { switchBotAPI = oldAPI }()

			ctx := context.Background()
			got, err := client.listDevice(ctx)

			if tt.wantErr {
				require.Error(t, err)
				if tt.wantErrMsg != "" {
					assert.Contains(t, err.Error(), tt.wantErrMsg)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				assert.Equal(t, tt.response.StatusCode, got.StatusCode)
				assert.Equal(t, tt.response.Message, got.Message)
				assert.Len(t, got.Body.DeviceList, len(tt.response.Body.DeviceList))
				assert.Len(t, got.Body.InfraredRemoteList, len(tt.response.Body.InfraredRemoteList))
			}
		})
	}
}
