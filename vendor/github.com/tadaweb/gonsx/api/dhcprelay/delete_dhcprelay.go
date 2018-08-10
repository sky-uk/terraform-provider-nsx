package dhcprelay

import (
	"github.com/tadaweb/gonsx/api"
	"net/http"
)

// DeleteDhcpRelayAPI - struct
type DeleteDhcpRelayAPI struct {
	*api.BaseAPI
}

// NewDelete - Generates a new DeleteDhcpRelayAPI object.
func NewDelete(edgeID string) *DeleteDhcpRelayAPI {
	this := new(DeleteDhcpRelayAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodDelete, "/api/4.0/edges/"+edgeID+"/dhcp/config/relay", nil, nil)
	return this
}
