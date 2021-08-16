package user

import (
	"encoding/json"
	"github.com/noc0lour/synapse_admin/api/types"
	"github.com/noc0lour/synapse_admin/pkg/client"
	"net/http"
	"net/url"
	// "fmt"
	// "io/ioutil"
)

type Client struct {
	*client.Client
}

func GetLastSeen(whois types.Whois) int {
	LastAccess := 0
	for _, dev := range whois.Devices {
		for _, sess := range dev.Sessions {
			for _, conn := range sess.Connections {
				if conn.LastSeen > LastAccess {
					LastAccess = conn.LastSeen
				}
			}
		}
	}
	return LastAccess
}

// Methods
func (c *Client) ListLastSeen(Since int) ([]string, error) {
	Since = Since * 1000 // Matrix uses Unix time in milliseconds
	var LastSeenUsers []string
	users, err := c.ListUsers()
	if err != nil {
		return LastSeenUsers, err
	}
	for _, u := range users {
		whois, err := c.WhoisUser(u.Name)
		if err != nil {
			return LastSeenUsers, err
		}
		if GetLastSeen(whois) >= Since {
			LastSeenUsers = append(LastSeenUsers, u.Name)
		}
	}
	return LastSeenUsers, nil

}
func (c *Client) WhoisUser(UserId string) (types.Whois, error) {
	var whois types.Whois
	// Build URL
	u, err := url.Parse(c.BuildBaseURL("_synapse", "admin", "v1", "whois", UserId))
	q := u.Query()
	q.Set("access_token", c.AccessToken)
	u.RawQuery = q.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return whois, err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := c.Client.Client.Client.Do(req)
	if err != nil {
		return whois, err
	}
	defer resp.Body.Close()

	// bodyBytes, err := ioutil.ReadAll(resp.Body)
	// bodyString := string(bodyBytes)
	// fmt.Printf("%v", bodyString)
	err = json.NewDecoder(resp.Body).Decode(&whois)
	return whois, err
}

func (c *Client) ListUsers() ([]types.User, error) {
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
	resp, err := c.Client.Client.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userlist types.UserList
	err = json.NewDecoder(resp.Body).Decode(&userlist)
	return userlist.Users, err
}

func (c *Client) DeactivateUser(UserId string) error {
	u, err := url.Parse(c.BuildBaseURL("_synapse", "admin", "v1", "deactivate", UserId))
	q := u.Query()
	q.Set("access_token", c.AccessToken)
	u.RawQuery = q.Encode()

	// Build Request
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return err
	}
	resp, err := c.Client.Client.Client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	return nil
}
