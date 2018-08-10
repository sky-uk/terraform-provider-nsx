package securitypolicy

import (
	"encoding/xml"
	"errors"
	"fmt"
)

func (sp SecurityPolicy) String() string {
	return fmt.Sprintf("SecurityPolicy with objectId: %s", sp.ObjectID)
}

// MarshalToXML converts the object into XML
func (sp SecurityPolicy) MarshalToXML() string {
	xmlBytes, _ := xml.Marshal(sp)
	return string(xmlBytes)
}

// AddSecurityGroupBinding - Adds security group to list of SecurityGroupBinding if it doesn't exists.
func (sp *SecurityPolicy) AddSecurityGroupBinding(objectID string) {
	for _, secGroup := range sp.SecurityGroupBinding {
		if secGroup.ObjectID == objectID {
			return
		}
	}
	// if we reached here that means we couldn't find one, and let's add the sec group.
	sp.SecurityGroupBinding = append(sp.SecurityGroupBinding, SecurityGroup{ObjectID: objectID})
	return
}

// RemoveSecurityGroupBinding - Adds security group to list of SecurityGroupBinding if it doesn't exists.
func (sp *SecurityPolicy) RemoveSecurityGroupBinding(objectID string) {
	for idx, secGroup := range sp.SecurityGroupBinding {
		if secGroup.ObjectID == objectID {
			sp.SecurityGroupBinding = append(sp.SecurityGroupBinding[:idx], sp.SecurityGroupBinding[idx+1:]...)
			return
		}
	}
	return
}

// CheckFirewallRuleByUUID - Checks if the rule with UUID exists in the firewall rules of security policy.
func (sp *SecurityPolicy) CheckFirewallRuleByUUID(uuid string) bool {
	for _, action := range sp.ActionsByCategory.Actions {
		if action.VsmUUID == uuid {
			return true
		}
	}
	return false
}

// GetFirewallRuleByName - Checks if the rule with given name exists in the firewall rules of security policy.
func (sp *SecurityPolicy) GetFirewallRuleByName(name string) *Action {
	var actionFound Action
	for _, action := range sp.ActionsByCategory.Actions {
		if action.Name == name {
			actionFound = action
			break
		}
	}
	return &actionFound
}

// GetFirewallRuleByUUID - Checks if the rule with given name exists in the firewall rules of security policy.
func (sp *SecurityPolicy) GetFirewallRuleByUUID(uuid string) *Action {
	var actionFound Action
	for _, action := range sp.ActionsByCategory.Actions {
		if action.VsmUUID == uuid {
			actionFound = action
			break
		}
	}
	return &actionFound
}

// RemoveFirewallActionByName - Removes the firewalla ction from security policy object if it exists.
func (sp *SecurityPolicy) RemoveFirewallActionByName(actionName string) {
	for idx, action := range sp.ActionsByCategory.Actions {
		if action.Name == actionName {
			sp.ActionsByCategory.Actions = append(sp.ActionsByCategory.Actions[:idx], sp.ActionsByCategory.Actions[idx+1:]...)
			return
		}
	}
}

// RemoveFirewallActionByUUID - Removes the firewall action from security policy object if it exists by it's UUID.
func (sp *SecurityPolicy) RemoveFirewallActionByUUID(uuid string) {
	for idx, action := range sp.ActionsByCategory.Actions {
		if action.VsmUUID == uuid {
			sp.ActionsByCategory.Actions = append(sp.ActionsByCategory.Actions[:idx], sp.ActionsByCategory.Actions[idx+1:]...)
			return
		}
	}
}

// AddFirewallAction adds inbount or outbound firewall action rule into security policy.
func (sp *SecurityPolicy) AddFirewallAction(name, action, direction string, secGroupObjectIDs, applicationObjectIDs []string) error {
	if action != "allow" && action != "block" && action != "reject" {
		return errors.New("Action can be only 'allow', 'block' or 'reject'")
	}
	if direction != "inbound" && direction != "outbound" {
		return errors.New("Direction can only be 'inbound' or 'outbound'")
	}

	var secondarySecurityGroupList = []SecurityGroup{}
	for _, secGroupID := range secGroupObjectIDs {
		securityGroup := SecurityGroup{ObjectID: secGroupID}
		secondarySecurityGroupList = append(secondarySecurityGroupList, securityGroup)
	}

	var secondaryApplicationsList = &Applications{}

	if applicationObjectIDs[0] != "any" {
		var secondaryApplicationList = []Application{}
		for _, applicationObjectID := range applicationObjectIDs {
			application := Application{ObjectID: applicationObjectID}
			secondaryApplicationList = append(secondaryApplicationList, application)
		}

		secondaryApplicationsList.Applications = secondaryApplicationList
	} else {
		secondaryApplicationsList = nil
	}

	newAction := Action{
		Class:                  "firewallSecurityAction",
		Name:                   name,
		Action:                 action,
		Category:               "firewall",
		Direction:              direction,
		IsEnabled:              true,
		SecondarySecurityGroup: secondarySecurityGroupList,
		NegateSource:           false,
		Applications:           secondaryApplicationsList,
	}

	if sp.ActionsByCategory.Category == "firewall" && len(sp.ActionsByCategory.Actions) > 0 {
		sp.ActionsByCategory.Actions = append(sp.ActionsByCategory.Actions, newAction)
		return nil
	}

	// Build actionsByCategory list.
	actionsByCategory := ActionsByCategory{Category: "firewall"}
	actionsByCategory.Actions = []Action{newAction}
	sp.ActionsByCategory = actionsByCategory
	return nil
}

// AddOutboundFirewallAction adds outbound firewall action rule into security policy.
// !! Deprecated in favor of AddFirewallAction
func (sp *SecurityPolicy) AddOutboundFirewallAction(name, action, direction string, secGroupObjectIDs, applicationObjectIDs []string) error {
	if action != "allow" && action != "block" {
		return errors.New("Action can be only 'allow' or 'block'")
	}
	if direction != "outbound" {
		return errors.New("Direction can only be 'outbound'")
	}

	var secondarySecurityGroupList = []SecurityGroup{}
	for _, secGroupID := range secGroupObjectIDs {
		securityGroup := SecurityGroup{ObjectID: secGroupID}
		secondarySecurityGroupList = append(secondarySecurityGroupList, securityGroup)
	}

	var secondaryApplicationsList = &Applications{}

	if applicationObjectIDs[0] != "any" {
		var secondaryApplicationList = []Application{}
		for _, applicationObjectID := range applicationObjectIDs {
			application := Application{ObjectID: applicationObjectID}
			secondaryApplicationList = append(secondaryApplicationList, application)
		}

		secondaryApplicationsList.Applications = secondaryApplicationList
	} else {
		secondaryApplicationsList = nil
	}

	newAction := Action{
		Class:                  "firewallSecurityAction",
		Name:                   name,
		Action:                 action,
		Category:               "firewall",
		Direction:              direction,
		IsEnabled:              true,
		SecondarySecurityGroup: secondarySecurityGroupList,
		Applications:           secondaryApplicationsList,
	}

	if sp.ActionsByCategory.Category == "firewall" && len(sp.ActionsByCategory.Actions) >= 1 {
		sp.ActionsByCategory.Actions = append(sp.ActionsByCategory.Actions, newAction)
		return nil
	}

	// Build actionsByCategory list.
	actionsByCategory := ActionsByCategory{Category: "firewall"}
	actionsByCategory.Actions = []Action{newAction}
	sp.ActionsByCategory = actionsByCategory
	return nil
}

// AddInboundFirewallAction adds outbound firewall action rule into security policy.
// !! Deprecated in favor of AddFirewallAction
func (sp *SecurityPolicy) AddInboundFirewallAction(name, action, direction string, applicationObjectIDs []string) error {
	if action != "allow" && action != "block" {
		return errors.New("Action can be only 'allow' or 'block'")
	}
	if direction != "inbound" {
		return errors.New("Direction can only be 'inbound'")
	}

	var secondaryApplicationsList = &Applications{}

	if applicationObjectIDs[0] != "any" {
		var secondaryApplicationList = []Application{}
		for _, applicationObjectID := range applicationObjectIDs {
			application := Application{ObjectID: applicationObjectID}
			secondaryApplicationList = append(secondaryApplicationList, application)
		}

		secondaryApplicationsList.Applications = secondaryApplicationList
	} else {
		secondaryApplicationsList = nil
	}

	newAction := Action{
		Class:        "firewallSecurityAction",
		Name:         name,
		Action:       action,
		Category:     "firewall",
		Direction:    direction,
		IsEnabled:    true,
		Applications: secondaryApplicationsList,
	}

	if sp.ActionsByCategory.Category == "firewall" && len(sp.ActionsByCategory.Actions) >= 1 {
		sp.ActionsByCategory.Actions = append(sp.ActionsByCategory.Actions, newAction)
		return nil
	}

	// Build actionsByCategory list.
	actionsByCategory := ActionsByCategory{Category: "firewall"}
	actionsByCategory.Actions = []Action{newAction}
	sp.ActionsByCategory = actionsByCategory
	return nil
}

func (spList SecurityPolicies) String() string {
	return fmt.Sprint("SecurityPolicies object, contains security policy objects.")
}

// FilterByName returns a single security policy object if it matches the name in SecurityPolicies list.
func (spList SecurityPolicies) FilterByName(name string) *SecurityPolicy {
	var securityPolicyFound SecurityPolicy
	for _, securityPolicy := range spList.SecurityPolicies {
		if securityPolicy.Name == name {
			securityPolicyFound = securityPolicy
			break
		}
	}
	return &securityPolicyFound
}

// RemoveSecurityPolicyByName - Removes the SecurityPolicy from a list of SecurityPolicies provided it matches the given name.
func (spList SecurityPolicies) RemoveSecurityPolicyByName(policyName string) *SecurityPolicies {
	for idx, securityPolicy := range spList.SecurityPolicies {
		if securityPolicy.Name == policyName {
			spList.SecurityPolicies = append(spList.SecurityPolicies[:idx], spList.SecurityPolicies[idx+1:]...)
			break
		}
	}
	return &spList
}
