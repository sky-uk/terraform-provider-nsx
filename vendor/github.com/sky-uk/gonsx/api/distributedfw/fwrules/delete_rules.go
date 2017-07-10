package fwrules

import (
	"github.com/sky-uk/gonsx/api"
	"net/http"
)

// DeleteFWRuleAPI default struct
type DeleteFWRuleAPI struct {
	*api.BaseAPI
}

// NewDelete - Returns all the rules in the specified context
func NewDelete(deleteRule Rule) *DeleteFWRuleAPI {
	this := new(DeleteFWRuleAPI)
	var endpoint string
	switch deleteRule.RuleType {
	case "LAYER3":
		endpoint = "/api/4.0/firewall/globalroot-0/config/layer3sections"

	case "LAYER2":
		endpoint = "/api/4.0/firewall/globalroot-0/config/layer2sections"
	}

	this.BaseAPI = api.NewBaseAPI(http.MethodDelete, endpoint, deleteRule, new(string))
	return this
}
