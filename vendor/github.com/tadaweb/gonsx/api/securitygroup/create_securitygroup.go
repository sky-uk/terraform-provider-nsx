package securitygroup

import (
	"github.com/tadaweb/gonsx/api"
	"net/http"
)

// CreateSecurityGroupAPI api object
type CreateSecurityGroupAPI struct {
	*api.BaseAPI
}

// NewCreate returns a new object of CreateSecurityGroupAPI.
func NewCreate(scopeID, securityGroupName string, dynamicMemberDefinition *DynamicMemberDefinition) *CreateSecurityGroupAPI {
	this := new(CreateSecurityGroupAPI)
	requestPayload := new(SecurityGroup)
	requestPayload.Name = securityGroupName

	requestPayload.DynamicMemberDefinition = dynamicMemberDefinition
	this.BaseAPI = api.NewBaseAPI(http.MethodPost, "/api/2.0/services/securitygroup/bulk/"+scopeID, requestPayload, new(string))
	return this
}

// GetResponse returns a ResponseObject of CreateSecurityGroupAPI.
func (ca CreateSecurityGroupAPI) GetResponse() string {
	return ca.ResponseObject().(string)
}
