package securitygroup

import (
	"github.com/tadaweb/gonsx/api"
	"net/http"
)

// UpdateSecurityGroupAPI object
type UpdateSecurityGroupAPI struct {
	*api.BaseAPI
}

// NewUpdate creates a new object of UpdateSecurityGroupAPI
func NewUpdate(securityGroupID string, securityGroupPayload *SecurityGroup) *UpdateSecurityGroupAPI {
	this := new(UpdateSecurityGroupAPI)
	endpointURL := "/api/2.0/services/securitygroup/bulk/" + securityGroupID
	this.BaseAPI = api.NewBaseAPI(http.MethodPut, endpointURL, securityGroupPayload, new(SecurityGroup))
	return this
}

// GetResponse returns the ResponseObject from UpdateServiceAPI
func (updateAPI UpdateSecurityGroupAPI) GetResponse() *SecurityGroup {
	return updateAPI.ResponseObject().(*SecurityGroup)
}
