package securitygroup

import (
	"github.com/tadaweb/gonsx/api"
	"net/http"
)

// GetAllSecurityGroupAPI base object.
type GetAllSecurityGroupAPI struct {
	*api.BaseAPI
}

// NewGetAll returns a new object of GetAllSecurityGroupAPI.
func NewGetAll(scopeID string) *GetAllSecurityGroupAPI {
	this := new(GetAllSecurityGroupAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodGet, "/api/2.0/services/securitygroup/scope/"+scopeID, nil, new(List))
	return this
}

// GetResponse returns ResponseObject of GetAllSecurityGroupAPI.
func (ga GetAllSecurityGroupAPI) GetResponse() *List {
	return ga.ResponseObject().(*List)
}
