package firewallexclusion

import (
	"github.com/tadaweb/gonsx/api"
	"net/http"
)

// GetAllFirewallExclusionAPI base object.
type GetAllFirewallExclusionAPI struct {
	*api.BaseAPI
}

// NewGetAll returns a new object of GetAllFirewallExclusionAPI.
func NewGetAll() *GetAllFirewallExclusionAPI {
	this := new(GetAllFirewallExclusionAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodGet, "/api/2.1/app/excludelist", nil, new(FirewallExclusions))
	return this
}

// GetResponse returns ResponseObject of GetAllIpSetAPI.
func (ga GetAllFirewallExclusionAPI) GetResponse() *FirewallExclusions {
	return ga.ResponseObject().(*FirewallExclusions)
}
