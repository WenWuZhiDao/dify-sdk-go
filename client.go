package dify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	host             string
	chatAPISecret    string
	datasetAPISecret string
	httpClient       *http.Client
}

func NewClientWithConfig(c *ClientConfig) *Client {
	var httpClient = &http.Client{}

	if c.Timeout != 0 {
		httpClient.Timeout = c.Timeout
	}
	if c.Transport != nil {
		httpClient.Transport = c.Transport
	}
	return &Client{
		host:             c.Host,
		chatAPISecret:    c.ChatAPISecret,
		datasetAPISecret: c.DatasetAPISecret,
		httpClient:       httpClient,
	}
}

func NewClient(host, chatAPISecret, datasetAPISecret string) *Client {
	return NewClientWithConfig(&ClientConfig{
		Host:             host,
		ChatAPISecret:    chatAPISecret,
		DatasetAPISecret: datasetAPISecret,
	})
}

func (c *Client) sendRequest(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}

func (c *Client) sendJSONRequest(req *http.Request, res interface{}) error {
	resp, err := c.sendRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP response StatusCode: %d, Body: %s ", resp.StatusCode, string(bodyBytes))
	}
	err = json.NewDecoder(resp.Body).Decode(res)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) getHost() string {
	var host = strings.TrimSuffix(c.host, "/")
	return host
}

func (c *Client) getChatAPISecret() string {
	return c.chatAPISecret
}

func (c *Client) getDatasetAPISecret() string {
	return c.datasetAPISecret
}

// Api deprecated, use API() instead
func (c *Client) Api() *API {
	return c.API()
}

func (c *Client) API() *API {
	return &API{
		c: c,
	}
}
