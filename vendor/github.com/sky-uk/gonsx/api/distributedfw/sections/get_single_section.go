package sections

import (
	"fmt"
	"github.com/sky-uk/gonsx/api"
	"net/http"
)

// GetSingleSectionAPI default struct
type GetSingleSectionAPI struct {
	*api.BaseAPI
}

// NewGetSingle - Returns all the rules in the specified context
func NewGetSingle(sectionID, sectionType string) *GetSingleSectionAPI {
	this := new(GetSingleSectionAPI)
	var endpoint string
	switch sectionType {
	case "LAYER3":
		endpoint = fmt.Sprintf("/api/4.0/firewall/globalroot-0/config/layer3sections/%s", sectionID)

	case "LAYER2":
		endpoint = fmt.Sprintf("/api/4.0/firewall/globalroot-0/config/layer2sections/%s", sectionID)
	}

	this.BaseAPI = api.NewBaseAPI(http.MethodGet, endpoint, nil, new(Section))

	return this

}

// GetSectionTimestamp - Returns the timestamp for a section , required to create a new rule
func GetSectionTimestamp(sectionID int, sectionType string) *GetSingleSectionAPI {
	this := new(GetSingleSectionAPI)
	var endpoint string
	switch sectionType {
	case "LAYER3":
		endpoint = fmt.Sprintf("/api/4.0/firewall/globalroot-0/config/layer3sections/%d", sectionID)
	case "LAYER2":
		endpoint = fmt.Sprintf("/api/4.0/firewall/globalroot-0/config/layer2sections/%d", sectionID)

	}
	this.BaseAPI = api.NewBaseAPI(http.MethodGet, endpoint, nil, new(Section))
	return this
}

// GetResponse - Returns ResponseObject from GetAllFirewallRulesAPI of Rule type.
func (getSingleAPI GetSingleSectionAPI) GetResponse() *Section {
	return getSingleAPI.ResponseObject().(*Section)
}
