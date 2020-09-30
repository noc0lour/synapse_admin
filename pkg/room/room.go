package room

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

func (c *Client) ListRooms() ([]types.Room, error) {
	// Build URL
	u, err := url.Parse(c.BuildBaseURL("_synapse", "admin", "v1", "rooms"))
	q := u.Query()
	q.Set("from", "0")
	// TODO implement pagination for big homeservers
	q.Set("limit", "99999999")
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

	var roomlist types.RoomList
	err = json.NewDecoder(resp.Body).Decode(&roomlist)
	return roomlist.Rooms, err
}
