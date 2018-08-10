package securitytag

import (
	"github.com/tadaweb/gonsx/api"
	"net/http"
)

// DetachSecurityTagAPI - struct
type DetachSecurityTagAPI struct {
	*api.BaseAPI
}

// NewDetach - Generates a new DetachSecurityTagAPI object.
func NewDetach(securityTagID, vmID string) *DetachSecurityTagAPI {
	this := new(DetachSecurityTagAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodDelete, "/api/2.0/services/securitytags/tag/"+securityTagID+"/vm/"+vmID, nil, nil)
	return this
}
