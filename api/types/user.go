package types

// https://github.com/matrix-org/synapse/blob/master/docs/admin_api/user_admin_api.rst
// Responses
type UserList struct {
	Users     []User `json:"users"`
	NextToken string `json:"next_token"`
	Total     int    `json:"total"`
}

type User struct {
	Name         string `json:"name"`
	PasswordHash string `json:"password_hash"`
	IsGuest      int    `json:"is_guest"`
	Admin        int    `json:"admin"`
	UserType     string `json:"user_type"`
	Deactivated  int    `json:"deactivated"`
	DisplayName  string `json:"displayname"`
	AvatarURL    string `json:"avatar_url"`
}

// https://github.com/matrix-org/synapse/blob/develop/docs/admin_api/user_admin_api.md

type Connection struct {
	Ip string `json:"ip"`
	LastSeen int `json:"last_seen"`
	UserAgent string `json:"user_agent"`
}

type Session struct {
	Connections []Connection `json:connections`
}

type Device struct {
	Sessions []Session  `json:sessions`
}

type Whois struct {
	UserId string `json:"user_id"`
	Devices map[string]Device `json:"devices"`
}
