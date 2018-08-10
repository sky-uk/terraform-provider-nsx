package securitypolicy

import (
	"github.com/tadaweb/gonsx/api"
	"net/http"
)

// DeleteSecurityPolicyAPI - struct
type DeleteSecurityPolicyAPI struct {
	*api.BaseAPI
}

// NewDelete - Generates a new DeleteSecurityPolicyAPI object.
func NewDelete(securityPolicyID string, force bool) *DeleteSecurityPolicyAPI {
	this := new(DeleteSecurityPolicyAPI)
	url := "/api/2.0/services/policy/securitypolicy/" + securityPolicyID + "?force=false"

	if force {
		url = "/api/2.0/services/policy/securitypolicy/" + securityPolicyID + "?force=true"
	}
	this.BaseAPI = api.NewBaseAPI(http.MethodDelete, url, nil, nil)
	return this
}
