package securitygroup

import (
	"github.com/sky-uk/gonsx/api"
	"net/http"
)

// CreateSecurityGroupAPI api object
type CreateSecurityGroupAPI struct {
	*api.BaseAPI
}

// NewCreate returns a new object of CreateServiceAPI.
func NewCreate(scopeID, securityGroupName string) *CreateSecurityGroupAPI {
	this := new(CreateSecurityGroupAPI)

	// Generate payload with the name.
	requestPayload := new(SecurityGroup)
	requestPayload.Name = securityGroupName

	// Build API object and return it
	this.BaseAPI = api.NewBaseAPI(http.MethodPost, "/api/2.0/services/securitygroup/bulk/"+scopeID, requestPayload, new(string))
	return this
}

// GetResponse returns a ResponseObject of CreateServiceAPI.
func (ca CreateSecurityGroupAPI) GetResponse() string {
	return ca.ResponseObject().(string)
}
