package sections

import (
	"fmt"
	"github.com/sky-uk/gonsx/api"
	"net/http"
)

// DeleteFWSectionAPI default struct
type DeleteFWSectionAPI struct {
	*api.BaseAPI
}

// NewDelete - Returns all the rules in the specified context
func NewDelete(deleteSection Section) *DeleteFWSectionAPI {
	this := new(DeleteFWSectionAPI)
	var endpoint string
	switch deleteSection.Type {
	case "LAYER3":
		endpoint = fmt.Sprintf("/api/4.0/firewall/globalroot-0/config/layer3sections/%s", deleteSection.ID)
	case "LAYER2":
		endpoint = fmt.Sprintf("/api/4.0/firewall/globalroot-0/config/layer2sections/%s", deleteSection.ID)
	}

	this.BaseAPI = api.NewBaseAPI(http.MethodDelete, endpoint, deleteSection, new(string))
	return this
}
