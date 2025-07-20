package api

import (
	"fmt"
	"io"
	"net/http"
)

type TokensClient client

// Validate validates the client's token
func (c *TokensClient) Validate() (bool, error) {
	r, err := c.client.Get("/v1/auth/validate", nil)
	if err != nil {
		return false, fmt.Errorf("failed to request validation: %s", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(r.Body)

	if r.StatusCode != http.StatusOK {
		return false, fmt.Errorf("failed to validate token: %w", parseResponseError(r))
	}

	data, err := unmarshal[struct{ Ok bool }](r)
	if err != nil {
		return false, fmt.Errorf("failed to deserialize validate token response: %w", err)
	}

	return data.Ok, nil
}

// Invalidate invalidates current token session
// TODO: @sanchitrk requires testing
func (c *TokensClient) Invalidate() (int64, error) {
	r, err := c.client.Post("/v1/auth/invalidate", nil)
	if err != nil {
		return 0, fmt.Errorf("failed to request invalidation: %s", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(r.Body)

	if r.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to invalidate sessions: %w", parseResponseError(r))
	}

	data, err := unmarshal[struct{ ValidFrom int64 }](r)
	if err != nil {
		return 0, fmt.Errorf("failed to deserialize invalidate sessions response: %w", err)
	}

	return data.ValidFrom, nil
}
