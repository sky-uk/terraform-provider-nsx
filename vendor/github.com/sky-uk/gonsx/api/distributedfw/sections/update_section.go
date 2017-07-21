package sections

import (
	"fmt"
	"github.com/sky-uk/gonsx/api"
	"net/http"
)

// UpdateFWSectionsAPI default struct
type UpdateFWSectionsAPI struct {
	*api.BaseAPI
}

// NewUpdate - Creates a new section
func NewUpdate(updateSection Section) *UpdateFWSectionsAPI {
	this := new(UpdateFWSectionsAPI)
	var endpoint string
	switch updateSection.Type {
	case "LAYER3":
		endpoint = fmt.Sprintf("/api/4.0/firewall/globalroot-0/config/layer3sections/%s", updateSection.ID)
	case "LAYER2":
		endpoint = fmt.Sprintf("/api/4.0/firewall/globalroot-0/config/layer2sections/%s", updateSection.ID)
	}
	this.BaseAPI = api.NewBaseAPI(http.MethodPut, endpoint, updateSection, new(Section))
	return this

}

// GetResponse - Returns ResponseObject from GetAllFirewallRulesAPI of Rule type.
func (updateAPI UpdateFWSectionsAPI) GetResponse() Section {
	return *updateAPI.ResponseObject().(*Section)
}
