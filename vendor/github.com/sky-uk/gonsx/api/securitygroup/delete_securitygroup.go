package securitygroup

import (
	"github.com/sky-uk/gonsx/api"
	"net/http"
)

// DeleteSecurityGroupAPI base object.
type DeleteSecurityGroupAPI struct {
	*api.BaseAPI
}

// NewDelete returns a new object of DeleteSecurityGroupAPI.
func NewDelete(serviceID string) *DeleteSecurityGroupAPI {
	this := new(DeleteSecurityGroupAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodDelete, "/api/2.0/services/securitygroup/"+serviceID, nil, nil)
	return this
}
