package user

import (
	"encoding/json"
	"github.com/noc0lour/synapse_admin/api/types"
	"github.com/noc0lour/synapse_admin/pkg/client"
	"net/http"
	"net/url"
)


type Client struct {
	*client.Client
}

// Methods
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

func (c *Client) DeactivateUser(UserId string) (error) {
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
