package fwrules

import (
	"fmt"
	"github.com/sky-uk/gonsx/api"
	"net/http"
)

// GetSingleRuleAPI default struct
type GetSingleRuleAPI struct {
	*api.BaseAPI
}

// NewGetSingle - Returns all the rules in the specified context
func NewGetSingle(ruleID, ruleType string, ruleSection int) *GetSingleRuleAPI {
	this := new(GetSingleRuleAPI)
	var endpoint string
	switch ruleType {
	case "LAYER3":
		endpoint = fmt.Sprintf("/api/4.0/firewall/globalroot-0/config/layer3sections/%d/rules/%s", ruleSection, ruleID)

	case "LAYER2":
		endpoint = fmt.Sprintf("/api/4.0/firewall/globalroot-0/config/layer2sections/%d/rules/%s", ruleSection, ruleID)
	}

	this.BaseAPI = api.NewBaseAPI(http.MethodGet, endpoint, nil, new(Rule))
	return this

}

// GetResponse - Returns ResponseObject from GetAllFirewallRulesAPI of Rule type.
func (getSingleAPI GetSingleRuleAPI) GetResponse() *Rule {
	return getSingleAPI.ResponseObject().(*Rule)
}
