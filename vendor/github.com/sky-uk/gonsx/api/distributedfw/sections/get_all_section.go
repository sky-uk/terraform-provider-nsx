package sections

import (
	"github.com/sky-uk/gonsx/api"
	"net/http"
)

// GetAllSectionsAPI default struct
type GetAllSectionsAPI struct {
	*api.BaseAPI
}

// NewGetAll - Returns all the rules in the specified context
func NewGetAll() *GetAllSectionsAPI {
	this := new(GetAllSectionsAPI)
	var endpoint string
	endpoint = "/api/4.0/firewall/globalroot-0/config"
	this.BaseAPI = api.NewBaseAPI(http.MethodGet, endpoint, nil, new(Section))
	return this

}

// GetResponse - Returns ResponseObject from GetAllFirewallRulesAPI of Rule type.
func (getAllAPI GetAllSectionsAPI) GetResponse() *Section {
	return getAllAPI.ResponseObject().(*Section)
}
