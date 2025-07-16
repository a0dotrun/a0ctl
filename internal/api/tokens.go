package api

import (
	"fmt"
	"io"
	"net/http"
)

type TokensClient client

// Validate validates the client's token
func (c *TokensClient) Validate() (int64, error) {
	r, err := c.client.Get("/v1/auth/validate", nil)
	if err != nil {
		return 0, fmt.Errorf("failed to request validation: %s", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(r.Body)

	if r.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to validate token: %w", parseResponseError(r))
	}

	data, err := unmarshal[struct{ Exp int64 }](r)
	if err != nil {
		return 0, fmt.Errorf("failed to deserialize validate token response: %w", err)
	}

	return data.Exp, nil
}
