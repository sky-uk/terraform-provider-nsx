package sections

import (
	"fmt"
	"github.com/sky-uk/gonsx/api"
	"net/http"
)

// CreateFWSectionsAPI default struct
type CreateFWSectionsAPI struct {
	*api.BaseAPI
}

// NewCreate - Creates a new section
func NewCreate(newSection Section) *CreateFWSectionsAPI {
	this := new(CreateFWSectionsAPI)
	var endpoint string
	switch newSection.Type {
	case "LAYER3":
		endpoint = fmt.Sprintf("/api/4.0/firewall/globalroot-0/config/layer3sections")
	case "LAYER2":
		endpoint = fmt.Sprintf("/api/4.0/firewall/globalroot-0/config/layer2sections")
	}
	this.BaseAPI = api.NewBaseAPI(http.MethodPost, endpoint, newSection, new(Section))
	return this

}

// GetResponse - Returns ResponseObject from GetAllFirewallRulesAPI of Rule type.
func (createAPI CreateFWSectionsAPI) GetResponse() Section {
	return *createAPI.ResponseObject().(*Section)
}
