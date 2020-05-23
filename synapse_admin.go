package main

import (
	"encoding/json"
	// "flag"
	"fmt"
	gomatrix "github.com/tulir/mautrix-go"
	"net/http"
	"net/url"
)


type Client struct {
	*gomatrix.Client
}

// https://github.com/matrix-org/synapse/blob/master/docs/admin_api/version_api.rst
func (c *Client) ServerVersion() error {
	url.Parse(c.BuildBaseURL("_synapse", "admin", "v1", "server_version"))
	return nil
}

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

// Methods
func (c *Client) ListUsers() ([]User, error) {
	// Build URL
	u, err := url.Parse(c.BuildBaseURL("_synapse", "admin", "v2", "users"))
	q := u.Query()
	q.Set("from", "0")
	// TODO implement pagination for big homeservers
	q.Set("limit", "99999999")
	q.Set("guests", "true")
	q.Set("access_token", c.AccessToken)
	u.RawQuery = q.Encode()

	// Build Request
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	// Execute Request
	resp, err := c.Client.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userlist UserList
	err = json.NewDecoder(resp.Body).Decode(&userlist)
	return userlist.Users, err
}

// https://github.com/matrix-org/synapse/blob/master/docs/admin_api/rooms.md
// Responses

type RoomList struct {
	Rooms []Room `json: "rooms"`
	Offset string `json "offset"`
	Total int `json "total_rooms"`


}
type Room struct {
	Id string `json: "room_id"`
	Name string `json "name"`
	Alias string `json "canonical_alias"`
	JoinedMembers int `json "joined_members"`
	JoinedLocalMembers int `json "joined_local_members"`
	Version string `json "version"`
	Creator string `json "creator"`
	Encryption string `json "encryption"`
	Federatable int `json "federatable"`
	Public int `json "public"`
	JoinRules string `json "join_rules"`
	GuestAccess string `json "GuestAccess"`
	HistoryVisibility string `json "HistoryVisibility"`
	StateEvents string `json "state_events"`

}
func (c *Client) ListRooms()([]Room, error) {
	// Build URL
	u, err := url.Parse(c.BuildBaseURL("_synapse", "admin", "v1", "rooms"))
	q := u.Query()
	q.Set("from", "0")
	// TODO implement pagination for big homeservers
	q.Set("limit", "99999999")
	q.Set("guests", "true")
	q.Set("access_token", c.AccessToken)
	u.RawQuery = q.Encode()


}

func main() {
	// Parse commandline arguments
	// Setup admin access
	cli, err := gomatrix.NewClient("https://matrix.example.com", "", "")
	user := gomatrix.UserIdentifier{Type: "m.id.user", User: "@user:example.com"}
	resp, err := cli.Login(&gomatrix.ReqLogin{Type: "m.login.password", Identifier: user, Password: "examplepassword"})
	if err != nil {
		fmt.Println(err)
		return
	}
	cli.SetCredentials(resp.UserID, resp.AccessToken)
	go cli.Sync()
	real_cli := Client{Client: cli}
	users, err := real_cli.ListUsers()
	if err != nil {
		fmt.Println("ListUsers() returned ", err)
		return
	}
	fmt.Printf("%+v", users)
}
