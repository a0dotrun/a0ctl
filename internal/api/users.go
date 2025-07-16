package api

import (
	"fmt"
	"net/http"
)

type UsersClient client

type UserInfo struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
}

func (c *UsersClient) GetUser() (UserInfo, error) {
	res, err := c.client.Get("/v1/user", nil)
	if err != nil {
		return UserInfo{}, fmt.Errorf("failed to get user info: %w", err)
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			return
		}
	}()

	if res.StatusCode != http.StatusOK {
		return UserInfo{}, parseResponseError(res)
	}

	data, err := unmarshal[UserInfo](res)
	if err != nil {
		return UserInfo{}, fmt.Errorf("failed to deserialize user response: %w", err)
	}

	return data, nil
}
