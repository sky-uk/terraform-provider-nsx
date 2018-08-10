package securitypolicy

import (
	"github.com/tadaweb/gonsx/api"
	"net/http"
)

// UpdateSecurityPolicyAPI ...
type UpdateSecurityPolicyAPI struct {
	*api.BaseAPI
}

// NewUpdate creates a new object of UpdateSecurityPolicyAPI
func NewUpdate(securityPolicyID string, securityPolicyPayload *SecurityPolicy) *UpdateSecurityPolicyAPI {
	this := new(UpdateSecurityPolicyAPI)
	endpointURL := "/api/2.0/services/policy/securitypolicy/" + securityPolicyID
	this.BaseAPI = api.NewBaseAPI(http.MethodPut, endpointURL, securityPolicyPayload, new(string))
	return this
}

// GetResponse returns the ResponseObject from UpdateSecurityPolicyAPI
func (updateAPI UpdateSecurityPolicyAPI) GetResponse() string {
	return updateAPI.ResponseObject().(string)
}
