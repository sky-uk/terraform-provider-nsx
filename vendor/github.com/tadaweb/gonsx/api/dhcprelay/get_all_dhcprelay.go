package dhcprelay

import (
	"github.com/tadaweb/gonsx/api"
	"net/http"
)

// GetAllDhcpRelaysAPI - struct
type GetAllDhcpRelaysAPI struct {
	*api.BaseAPI
}

// NewGetAll  - returns GetAll api object of GetAllDhcpRelaysAPI type.
func NewGetAll(edgeID string) *GetAllDhcpRelaysAPI {
	this := new(GetAllDhcpRelaysAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodGet, "/api/4.0/edges/"+edgeID+"/dhcp/config/relay", nil, new(DhcpRelay))
	return this
}

// GetResponse - Returns ResponseObject from GetAllDhcpRelaysAPI of DhcpRelay type.
func (getAllAPI GetAllDhcpRelaysAPI) GetResponse() *DhcpRelay {
	return getAllAPI.ResponseObject().(*DhcpRelay)
}
