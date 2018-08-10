package firewallexclusion

import (
	"github.com/tadaweb/gonsx/api"
	"net/http"
)

// CreateFirewallExclusionAPI base object.
type CreateFirewallExclusionAPI struct {
	*api.BaseAPI
}

// NewCreate returns a new object of CreateFirewallExclusionAPI.
func NewCreate(moid string) *CreateFirewallExclusionAPI {
	this := new(CreateFirewallExclusionAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodPut, "/api/2.1/app/excludelist/"+moid, nil, nil)
	return this
}

// GetResponse returns ResponseObject of CreateFirewallExclusionAPI.
func (ga CreateFirewallExclusionAPI) GetResponse() string {
	return ga.ResponseObject().(string)
}
