package securitypolicy

import (
	"github.com/tadaweb/gonsx/api"
	"net/http"
)

// GetSecurityPolicyAPI base object.
type GetSecurityPolicyAPI struct {
	*api.BaseAPI
}

// NewGet returns new object of GetSecurityPolicyAPI.
func NewGet(securityPolicyID string) *GetSecurityPolicyAPI {
	this := new(GetSecurityPolicyAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodGet, "/api/2.0/services/policy/securitypolicy/"+securityPolicyID, nil, new(SecurityPolicy))
	return this
}

// GetResponse returns ResponseObject of GetSecurityPolicyAPI.
func (ga GetSecurityPolicyAPI) GetResponse() *SecurityPolicy {
	return ga.ResponseObject().(*SecurityPolicy)
}
