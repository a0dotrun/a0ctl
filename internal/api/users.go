package api

type UsersClient client

type UserInfo struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
}
