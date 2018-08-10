package securitytag

import (
	"github.com/tadaweb/gonsx/api"
	"net/http"
)

// CreateSecurityTagAPI - struct
type CreateSecurityTagAPI struct {
	*api.BaseAPI
}

// NewCreate - Generates a new CreateSecurityTagAPI object.
func NewCreate(name, desc string) *CreateSecurityTagAPI {
	this := new(CreateSecurityTagAPI)
	requestPayload := new(SecurityTag)
	requestPayload.Name = name
	requestPayload.Description = desc
	// TODO: need to make it argument
	requestPayload.TypeName = "SecurityTag"
	this.BaseAPI = api.NewBaseAPI(http.MethodPost, "/api/2.0/services/securitytags/tag", requestPayload, new(string))
	return this
}

// GetResponse returns the ResponseObject from CreateSecurityTagAPI
func (createAPI CreateSecurityTagAPI) GetResponse() string {
	return createAPI.ResponseObject().(string)
}
