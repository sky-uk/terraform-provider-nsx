package dhcprelay

import (
	"github.com/sky-uk/gonsx/api"
	"net/http"
)

// UpdateDhcpRelayAPI ...
type UpdateDhcpRelayAPI struct {
	*api.BaseAPI
}

// NewUpdate creates a new object of UpdateDhcpRelayAPI
func NewUpdate(dhcpIPAddress, edgeID string, relayAgentslist []RelayAgent) *UpdateDhcpRelayAPI {
	this := new(UpdateDhcpRelayAPI)
	requestPayload := new(DhcpRelay)
	requestPayload.RelayServer.IPAddress = dhcpIPAddress
	requestPayload.RelayAgents = relayAgentslist

	this.BaseAPI = api.NewBaseAPI(http.MethodPut, "/api/4.0/edges/"+edgeID+"/dhcp/config/relay", requestPayload, new(string))
	return this
}

// GetResponse returns the ResponseObject from UpdateDhcpRelayAPI
func (updateAPI UpdateDhcpRelayAPI) GetResponse() string {
	return updateAPI.ResponseObject().(string)
}
