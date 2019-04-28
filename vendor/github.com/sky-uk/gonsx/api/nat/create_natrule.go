package nat

import (
	"github.com/sky-uk/gonsx/api"
	"net/http"
)

// CreateNatRuleAPI ...
type CreateNatRuleAPI struct {
	*api.BaseAPI
}

// NewCreateRule creates a new object of UpdateDhcpRelayAPI
func NewCreateRule(edgeID string, natRule Rule) *CreateNatRuleAPI {
	this := new(CreateNatRuleAPI)
	var requestPayload Rules
	requestPayload.Rules = append(requestPayload.Rules, natRule)

	this.BaseAPI = api.NewBaseAPI(http.MethodPost, "/api/4.0/edges/"+edgeID+"/nat/config/rules", requestPayload, new(string))
	return this
}

// NewCreateRules creates a new object of UpdateDhcpRelayAPI
func NewCreateRules(edgeID string, natRules Rules) *CreateNatRuleAPI {
	this := new(CreateNatRuleAPI)
	requestPayload := natRules
	this.BaseAPI = api.NewBaseAPI(http.MethodPost, "/api/4.0/edges/"+edgeID+"/nat/config/rules", requestPayload, new(string))
	return this
}

// GetResponse returns the ResponseObject from UpdateDhcpRelayAPI
func (createAPI CreateNatRuleAPI) GetResponse() string {
	return createAPI.ResponseObject().(string)
}
