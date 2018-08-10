package securitygroup

import (
	"github.com/tadaweb/gonsx/api"
	"net/http"
)

// DeleteSecurityGroupAPI base object.
type DeleteSecurityGroupAPI struct {
	*api.BaseAPI
}

// NewDelete returns a new object of DeleteSecurityGroupAPI.
func NewDelete(securityGroupID string) *DeleteSecurityGroupAPI {
	this := new(DeleteSecurityGroupAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodDelete, "/api/2.0/services/securitygroup/"+securityGroupID, nil, nil)
	return this
}
