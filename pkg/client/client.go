package client
import (
	gomatrix "maunium.net/go/mautrix"
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
