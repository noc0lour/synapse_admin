package types

// https://github.com/matrix-org/synapse/blob/master/docs/admin_api/user_admin_api.rst
// Responses
type UserList struct {
	Users     []User `json: "users"`
	NextToken string `json: "next_token"`
	Total     int    `json: "total"`
}

type User struct {
	Name         string `json: "name"`
	PasswordHash string `json: "password_hash"`
	IsGuest      int    `json: "is_guest"`
	Admin        int    `json: "admin"`
	UserType     string `json: "user_type"`
	Deactivated  int    `json: "deactivated"`
	DisplayName  string `json: "displayname"`
	AvatarURL    string `json: "avatar_url"`
}
