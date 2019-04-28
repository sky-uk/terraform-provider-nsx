package nat

import (
	"github.com/sky-uk/gonsx/api"
	"net/http"
)

// DeleteNatRuleAPI - struct
type DeleteNatRuleAPI struct {
	*api.BaseAPI
}

// NewDelete - Generates a new DeleteNatRuleAPI object.
func NewDelete(edgeID string, natRuleID string) *DeleteNatRuleAPI {
	this := new(DeleteNatRuleAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodDelete, "/api/4.0/edges/"+edgeID+"/nat/config/rules/"+natRuleID, nil, nil)
	return this
}
