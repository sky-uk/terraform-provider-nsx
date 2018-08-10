package dhcprelay

import (
	"github.com/tadaweb/gonsx/api"
	"net/http"
)

// CreateDhcpRelayAPI ...
type CreateDhcpRelayAPI struct {
	*api.BaseAPI
}

// NewCreate creates a new object of UpdateDhcpRelayAPI
func NewCreate(edgeID string, dhcpRelay DhcpRelay) *CreateDhcpRelayAPI {
	this := new(CreateDhcpRelayAPI)
	requestPayload := dhcpRelay
	this.BaseAPI = api.NewBaseAPI(http.MethodPut, "/api/4.0/edges/"+edgeID+"/dhcp/config/relay", requestPayload, new(string))
	return this
}

// GetResponse returns the ResponseObject from UpdateDhcpRelayAPI
func (createAPI CreateDhcpRelayAPI) GetResponse() string {
	return createAPI.ResponseObject().(string)
}
