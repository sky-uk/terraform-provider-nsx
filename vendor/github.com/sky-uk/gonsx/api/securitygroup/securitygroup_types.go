package securitygroup

import "encoding/xml"

// List - top level <list> element
type List struct {
	XMLName        xml.Name        `xml:"list"`
	SecurityGroups []SecurityGroup `xml:"securitygroup"`
}

// SecurityGroup - <securitygroup> element of <list>
type SecurityGroup struct {
	XMLName                 xml.Name                `xml:"securitygroup"`
	ObjectID                string                  `xml:"objectId,omitempty"`
	ObjectTypeName          string                  `xml:"objectTypeName,omitempty"`
	Revision                string                  `xml:"revision,omitempty"`
	Type                    string                  `xml:"type,omitempty>typeName,omitempty"`
	Name                    string                  `xml:"name"`
	InheritanceAllowed      bool                    `xml:"inheritanceAllowed,omitempty"`
	DynamicMemberDefinition DynamicMemberDefinition `xml:"dynamicMemberDefinition"`
}

// DynamicMemberDefinition - <dynamicMemberDefinition> element of <securitygroup>
type DynamicMemberDefinition struct {
	DynamicSet []DynamicSet `xml:"dynamicSet"`
}

// DynamicSet - <dynamicSet> element of <dynamicMemberDefinition>
type DynamicSet struct {
	Operator        string            `xml:"operator"`
	DynamicCriteria []DynamicCriteria `xml:"dynamicCriteria"`
}

// DynamicCriteria - <dynamicCriteria> element of <dynamicSet>
type DynamicCriteria struct {
	Operator string `xml:"operator"`
	Key      string `xml:"key"`
	Criteria string `xml:"criteria"`
	Value    string `xml:"value"`
	IsValid  bool   `xml:"isValid,omitempty"`
}
