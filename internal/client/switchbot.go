package client

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	switchBotAPI = "https://api.switch-bot.com/v1.1"
)

type Client struct {
	HttpClient *http.Client

	token  string
	secret string
}

func NewClient(token, secret string) *Client {
	return &Client{
		HttpClient: &http.Client{},
		token:      token,
		secret:     secret,
	}
}

func (c *Client) get(ctx context.Context, path string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, switchBotAPI+path, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	return c.Do(req)
}

func (c *Client) post(ctx context.Context, path string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, switchBotAPI+path, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	return c.Do(req)
}

func (c *Client) Do(req *http.Request) ([]byte, error) {
	nonce := uuid.New().String()
	timestamp := time.Now().UnixMilli()
	data := fmt.Sprintf("%s%d%s", c.token, timestamp, nonce)
	mac := hmac.New(sha256.New, []byte(c.secret))
	mac.Write([]byte(data))
	signature := mac.Sum(nil)
	signatureB64 := strings.ToUpper(base64.StdEncoding.EncodeToString(signature))

	req.Header.Set("Authorization", c.token)
	req.Header.Set("sign", signatureB64)
	req.Header.Set("nonce", nonce)
	req.Header.Set("t", fmt.Sprintf("%d", timestamp))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			// Log the error but don't return it as the main operation may have succeeded
			// In production, you might want to use a proper logger here
			_ = err
		}
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}
	return body, nil
}
