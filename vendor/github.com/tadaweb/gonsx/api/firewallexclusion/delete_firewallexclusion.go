package firewallexclusion

import (
	"github.com/tadaweb/gonsx/api"
	"net/http"
)

// DeleteFirewallExclusionAPI base object.
type DeleteFirewallExclusionAPI struct {
	*api.BaseAPI
}

// NewDelete returns a new object of DeleteFirewallExclusionAPI.
func NewDelete(moid string) *DeleteFirewallExclusionAPI {
	this := new(DeleteFirewallExclusionAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodDelete, "/api/2.1/app/excludelist/"+moid, nil, nil)
	return this
}

// GetResponse returns ResponseObject of DeleteFirewallExclusionAPI.
func (ga DeleteFirewallExclusionAPI) GetResponse() string {
	return ga.ResponseObject().(string)
}
