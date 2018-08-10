package securitytag

import (
	"github.com/tadaweb/gonsx/api"
	"net/http"
)

// DeleteSecurityTagAPI - struct
type DeleteSecurityTagAPI struct {
	*api.BaseAPI
}

// NewDelete - Generates a new DeleteSecurityTagApi object.
func NewDelete(securityTagID string) *DeleteSecurityTagAPI {
	this := new(DeleteSecurityTagAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodDelete, "/api/2.0/services/securitytags/tag/"+securityTagID, nil, nil)
	return this
}
