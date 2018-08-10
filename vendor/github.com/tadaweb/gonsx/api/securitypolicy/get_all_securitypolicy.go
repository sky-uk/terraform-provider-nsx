package securitypolicy

import (
	"github.com/tadaweb/gonsx/api"
	"net/http"
)

// GetAllSecurityPoliciesAPI - struct
type GetAllSecurityPoliciesAPI struct {
	*api.BaseAPI
}

// NewGetAll  - returns GetAll api object of GetAllSecurityPoliciesAPI type.
func NewGetAll() *GetAllSecurityPoliciesAPI {
	this := new(GetAllSecurityPoliciesAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodGet, "/api/2.0/services/policy/securitypolicy/all", nil, new(SecurityPolicies))
	return this
}

// GetResponse - Returns ResponseObject from GetAllSecurityPoliciesAPI.
func (getAllAPI GetAllSecurityPoliciesAPI) GetResponse() *SecurityPolicies {
	return getAllAPI.ResponseObject().(*SecurityPolicies)
}
