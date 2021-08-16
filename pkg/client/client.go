package client

import (
	"github.com/adrg/xdg"
	gomatrix "maunium.net/go/mautrix"
	"net/url"
	"path"
)

type Client struct {
	*gomatrix.Client
}

type credentials struct {
	Server      string
	UserID      string
	AccessToken string
}

// https://github.com/matrix-org/synapse/blob/master/docs/admin_api/version_api.rst
func (c *Client) ServerVersion() error {
	url.Parse(c.BuildBaseURL("_synapse", "admin", "v1", "server_version"))
	return nil
}

func (c *Client) SaveCredentials() error {
	cache_path := path.Join(xdg.CacheHome, "synapse_admin")
	_ = cache_path

	return nil
}

func (c *Client) LoadCredentials() error {

	return nil
}
