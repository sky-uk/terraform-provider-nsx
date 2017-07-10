package fwrules

import (
	"fmt"
	"github.com/sky-uk/gonsx/api"
	"net/http"
)

// CreateFWRulesAPI default struct
type CreateFWRulesAPI struct {
	*api.BaseAPI
}

// NewCreate - Returns all the rules in the specified context
func NewCreate(newRule Rule) *CreateFWRulesAPI {
	this := new(CreateFWRulesAPI)
	var endpoint string
	switch newRule.RuleType {
	case "LAYER3":
		endpoint = fmt.Sprintf("/api/4.0/firewall/globalroot-0/config/layer3sections/%d/rules", newRule.SectionID)

	case "LAYER2":
		endpoint = fmt.Sprintf("/api/4.0/firewall/globalroot-0/config/layer2sections/%d/rules", newRule.SectionID)

	}

	this.BaseAPI = api.NewBaseAPI(http.MethodPost, endpoint, newRule, new(Rule))
	return this
}

// GetResponse - Returns ResponseObject from GetAllFirewallRulesAPI of Rule type.
func (createAPI CreateFWRulesAPI) GetResponse() Rule {
	return *createAPI.ResponseObject().(*Rule)
}
