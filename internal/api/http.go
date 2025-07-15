package api

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func Header(key, value string) map[string]string {
	return map[string]string{
		key: value,
	}
}

func (c *Client) Get(path string, body io.Reader) (*http.Response, error) {
	return c.do("GET", path, body, Header("Content-Type", "application/json"))
}

func (c *Client) GetWithHeaders(
	path string, body io.Reader, headers map[string]string,
) (*http.Response, error) {
	headers["Content-Type"] = "application/json"
	return c.do("GET", path, body, headers)
}

func (c *Client) Post(path string, body io.Reader) (*http.Response, error) {
	return c.do("POST", path, body, Header("Content-Type", "application/json"))
}

func (c *Client) PostBinary(path string, body io.Reader) (*http.Response, error) {
	return c.do("POST", path, body, Header("Content-Type", "application/octet-stream"))
}

func (c *Client) Patch(path string, body io.Reader) (*http.Response, error) {
	return c.do("PATCH", path, body, Header("Content-Type", "application/json"))
}

func (c *Client) Put(path string, body io.Reader) (*http.Response, error) {
	return c.do("PUT", path, body, Header("Content-Type", "application/json"))
}

func (c *Client) Upload(path string, fileData *os.File) (*http.Response, error) {
	body, bodyWriter := io.Pipe()
	writer := multipart.NewWriter(bodyWriter)
	go func() {
		formFile, err := writer.CreateFormFile("file", fileData.Name())
		if err != nil {
			bodyWriter.CloseWithError(err)
			return
		}
		if _, err := io.Copy(formFile, fileData); err != nil {
			bodyWriter.CloseWithError(err)
			return
		}
		bodyWriter.CloseWithError(writer.Close())
	}()
	req, err := c.newRequest("POST", path, body, Header("Content-Type", writer.FormDataContentType()))
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) Delete(path string, body io.Reader) (*http.Response, error) {
	return c.do("DELETE", path, body, Header("Content-Type", "application/json"))
}
