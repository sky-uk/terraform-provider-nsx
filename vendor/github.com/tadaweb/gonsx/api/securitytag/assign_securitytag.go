package securitytag

import (
	"github.com/tadaweb/gonsx/api"
	"net/http"
)

// AssignSecurityTagAPI struct
type AssignSecurityTagAPI struct {
	*api.BaseAPI
}

// NewAssign - Generates a new AssignSecurityTagAPI object.
func NewAssign(securityTagID, vmID string) *AssignSecurityTagAPI {
	this := new(AssignSecurityTagAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodPut, "/api/2.0/services/securitytags/tag/"+securityTagID+"/vm/"+vmID, nil, nil)
	return this
}
