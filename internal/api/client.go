// Package api provides the API client definition and implementations.
package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"runtime"

	"github.com/a0dotrun/a0ctl/internal/config"

	"github.com/a0dotrun/a0ctl/internal/flags"
)

type ErrorResponseDetails struct {
	Error interface{} `json:"error"`
	Code  string      `json:"code"`
}

func unmarshal[T any](r *http.Response) (T, error) {
	d, err := io.ReadAll(r.Body)
	t := new(T)
	if err != nil {
		return *t, err
	}
	err = json.Unmarshal(d, &t)
	return *t, err
}

func marshal(data interface{}) (io.Reader, error) {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(data)
	return buf, err
}

// Client represents the API client for a0ctl.
type Client struct {
	BaseURL    *url.URL
	Token      string
	Username   string
	CLIVersion string

	// Single instance to be reused by all clients
	base *client

	Tokens *TokensClient
}

// client struct that will be aliases by all other clients
type client struct {
	client *Client
}

func NewClient(baseURL *url.URL, token, username string) *Client {
	c := &Client{
		BaseURL:    baseURL,
		Token:      token,
		Username:   username,
		CLIVersion: "0.0.1", // FIXME: read from config
	}

	c.base = &client{client: c}
	c.Tokens = (*TokensClient)(c.base)

	return c
}

// AuthedClient returns authenticated client
func AuthedClient() (*Client, error) {
	token, err := GetAccessToken()
	if err != nil {
		return nil, err
	}
	return MakeClient(token)
}

// MakeClient builds a new API client with the provided token.
// Reads settings from configdir.
func MakeClient(token string) (*Client, error) {
	urlStr := config.GetA0URL()
	a0URL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("error creating a0ctl client: could not parse a0 URL %s: %w", urlStr, err)
	}

	settings, err := config.ReadSettings()
	if err != nil {
		return nil, fmt.Errorf("error creating a0ctl client: could not read settings: %w", err)
	}

	username := settings.GetUsername()
	return NewClient(a0URL, token, username), nil
}

func (c *Client) newRequest(
	method, urlPath string, body io.Reader, extraHeaders map[string]string,
) (*http.Request, error) {
	if _, exists := extraHeaders["Content-Type"]; !exists {
		return nil, errors.New("content type is required")
	}
	baseURL, err := url.Parse(c.BaseURL.String())
	if err != nil {
		return nil, err
	}
	baseURL, err = baseURL.Parse(urlPath)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, baseURL.String(), body)
	if err != nil {
		return nil, err
	}
	if c.Token != "" {
		req.Header.Add("Authorization", fmt.Sprint("Bearer ", c.Token))
	}
	req.Header.Add("a0ctlversion", c.CLIVersion)

	parsedCliVersion := c.CLIVersion
	if parsedCliVersion != "dev" {
		parsedCliVersion = c.CLIVersion[1:] // strip the leading "v"
	}
	req.Header.Add(
		"User-Agent",
		fmt.Sprintf("a0ctl/%s (%s/%s)",
			parsedCliVersion, runtime.GOOS, runtime.GOARCH),
	)
	for header, value := range extraHeaders {
		req.Header.Add(header, value)
	}
	return req, nil
}

func (c *Client) do(
	method, path string, body io.Reader, extraHeaders map[string]string,
) (*http.Response, error) {
	req, err := c.newRequest(method, path, body, extraHeaders)
	var reqDump string
	if flags.Debug() {
		reqDump = dumpRequest(req)
	}
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if flags.Debug() {
		printDumps(reqDump, dumpResponse(resp))
	}
	return resp, nil
}

func dumpRequest(req *http.Request) string {
	dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return ""
	}
	return string(dump)
}

func dumpResponse(req *http.Response) string {
	dump, err := httputil.DumpResponse(req, true)
	if err != nil {
		return ""
	}
	return string(dump)
}

func printDumps(req, resp string) {
	if req != "" {
		fmt.Println(req)
	}
	if resp != "" {
		fmt.Println(resp)
	}
}

func parseResponseError(res *http.Response) error {
	d, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("response failed with status %s", res.Status)
	}

	var errResp ErrorResponseDetails
	if err := json.Unmarshal(d, &errResp); err == nil {
		if errResp.Error != nil {
			return fmt.Errorf("%v", errResp.Error)
		}
	}
	return fmt.Errorf("response failed with status %s", res.Status)
}
