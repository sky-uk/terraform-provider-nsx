package securitytag

import (
	"github.com/tadaweb/gonsx/api"
	"net/http"
)

// GetAllAttachedSecurityTagsAPI - struct
type GetAllAttachedSecurityTagsAPI struct {
	*api.BaseAPI
}

// NewGetAllAttached - Generates a new GetAllSecurityTagsAttachedAPI object.
func NewGetAllAttached(tagID string) *GetAllAttachedSecurityTagsAPI {
	this := new(GetAllAttachedSecurityTagsAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodGet, "/api/2.0/services/securitytags/tag/"+tagID+"/vm", nil, new(BasicInfoList))
	return this
}

// GetResponse returns the ResponseObject from GetAllAttachedSecurityTagsAPI
func (getAPI GetAllAttachedSecurityTagsAPI) GetResponse() *BasicInfoList {
	return getAPI.ResponseObject().(*BasicInfoList)
}
