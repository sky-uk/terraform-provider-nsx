package dhcprelay

import (
	"github.com/tadaweb/gonsx/api"
	"net/http"
)

// UpdateDhcpRelayAPI ...
type UpdateDhcpRelayAPI struct {
	*api.BaseAPI
}

// NewUpdate creates a new object of UpdateDhcpRelayAPI
func NewUpdate(edgeID string, dhcpRelay DhcpRelay) *UpdateDhcpRelayAPI {
	this := new(UpdateDhcpRelayAPI)
	requestPayload := dhcpRelay
	this.BaseAPI = api.NewBaseAPI(http.MethodPut, "/api/4.0/edges/"+edgeID+"/dhcp/config/relay", requestPayload, new(string))
	return this
}

// GetResponse returns the ResponseObject from UpdateDhcpRelayAPI
func (updateAPI UpdateDhcpRelayAPI) GetResponse() string {
	return updateAPI.ResponseObject().(string)
}
