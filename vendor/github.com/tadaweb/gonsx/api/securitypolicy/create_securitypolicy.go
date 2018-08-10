package securitypolicy

import (
	"github.com/tadaweb/gonsx/api"
	"net/http"
)

// CreateSecurityPolicyAPI api object
type CreateSecurityPolicyAPI struct {
	*api.BaseAPI
}

// NewCreate returns a new object of CreatePolicyAPI.
func NewCreate(name, precedence, description string, securityGroups []string, actions []Action) *CreateSecurityPolicyAPI {
	this := new(CreateSecurityPolicyAPI)
	requestPayload := new(SecurityPolicy)
	requestPayload.Name = name
	requestPayload.Precedence = precedence
	requestPayload.Description = description
	requestPayload.SecurityGroupBinding = []SecurityGroup{}

	var securityGroupBindingList = []SecurityGroup{}
	for _, secGroupID := range securityGroups {
		securityGroupBinding := SecurityGroup{ObjectID: secGroupID}
		securityGroupBindingList = append(securityGroupBindingList, securityGroupBinding)
	}
	requestPayload.SecurityGroupBinding = securityGroupBindingList
	requestPayload.ActionsByCategory = ActionsByCategory{Actions: actions}
	this.BaseAPI = api.NewBaseAPI(http.MethodPost, "/api/2.0/services/policy/securitypolicy", requestPayload, new(string))
	return this
}

// GetResponse returns a ResponseObject of CreateServiceAPI.
func (ca CreateSecurityPolicyAPI) GetResponse() string {
	return ca.ResponseObject().(string)
}
