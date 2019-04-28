package nat

import (
	"github.com/sky-uk/gonsx/api"
	"net/http"
)

// UpdateNatRuleAPI ...
type UpdateNatRuleAPI struct {
	*api.BaseAPI
}

// NewUpdate creates a new object of UpdateNatRuleAPI
func NewUpdate(edgeID string, natRuleID string, natRule Rule) *UpdateNatRuleAPI {
	this := new(UpdateNatRuleAPI)
	requestPayload := natRule
	this.BaseAPI = api.NewBaseAPI(http.MethodPut, "/api/4.0/edges/"+edgeID+"/nat/config/rules/"+natRuleID, requestPayload, new(string))
	return this
}

// GetResponse returns the ResponseObject from UpdateNatRuleAPI
func (updateAPI UpdateNatRuleAPI) GetResponse() string {
	return updateAPI.ResponseObject().(string)
}
