package securitypolicy

import "encoding/xml"

// SecurityPolicies top level struct
type SecurityPolicies struct {
	SecurityPolicies []SecurityPolicy `xml:"securityPolicy"`
}

// SecurityPolicy object struct
type SecurityPolicy struct {
	XMLName              xml.Name          `xml:"securityPolicy"`
	ObjectID             string            `xml:"objectId,omitempty"`
	ObjectTypeName       string            `xml:"objectTypeName,omitempty"`
	VsmUUID              string            `xml:"vsmUuid,omitempty"`
	NodeID               string            `xml:"nodeId,omitempty"`
	Revision             int               `xml:"revision,omitempty"`
	TypeName             string            `xml:"type,omitempty>typeName,omitempty"`
	Name                 string            `xml:"name,omitempty"`
	Description          string            `xml:"description,omitempty"`
	Precedence           string            `xml:"precedence"`
	IsUniversal          bool              `xml:"isUniversal,omitempty"`
	InheritanceAllowed   bool              `xml:"inheritanceAllowed,omitempty"`
	ActionsByCategory    ActionsByCategory `xml:"actionsByCategory,omitempty"`
	SecurityGroupBinding []SecurityGroup   `xml:"securityGroupBinding,omitempty"`
}

// SecurityGroup object struct
type SecurityGroup struct {
	ObjectID string `xml:"objectId,omitempty"`
}

// ActionsByCategory element of SecurityPolicy.
type ActionsByCategory struct {
	XMLName  xml.Name `xml:"actionsByCategory"`
	Category string   `xml:"category,omitempty"`
	Actions  []Action `xml:"action,omitempty"`
}

// Action element of ActionsByCategory list.
type Action struct {
	XMLName                xml.Name        `xml:"action"`
	Class                  string          `xml:"class,attr"`
	ObjectID               string          `xml:"objectId,omitempty"`
	ObjectTypeName         string          `xml:"objectTypeName,omitempty"`
	VsmUUID                string          `xml:"vsmUuid,omitempty"`
	NodeID                 string          `xml:"nodeId,omitempty"`
	Revision               int             `xml:"revision,omitempty"`
	TypeName               string          `xml:"type,omitempty>typeName,omitempty"`
	Name                   string          `xml:"name,omitempty"`
	Action                 string          `xml:"action"`
	Category               string          `xml:"category"`
	Direction              string          `xml:"direction"`
	IsEnabled              bool            `xml:"isEnabled,omitempty"`
	SecondarySecurityGroup []SecurityGroup `xml:"secondarySecurityGroup,omitempty"`
	NegateSource           bool            `xml:"outsideSecondaryContainer,omitempty"`
	Applications           *Applications   `xml:"applications,omitempty"`
}

//Applications object struct
type Applications struct {
	Applications []Application `xml:"application,omitempty"`
}

// Application element of Applications List
type Application struct {
	//XMLName  xml.Name `xml:"application,omitempty"`
	ObjectID string `xml:"objectId,omitempty"`
	Name     string `xml:"name,omitempty"`
}
