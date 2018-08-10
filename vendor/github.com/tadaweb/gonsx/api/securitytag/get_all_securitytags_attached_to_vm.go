package securitytag

import (
	"github.com/tadaweb/gonsx/api"
	"net/http"
)

// GetAllSecurityTagsAttachedToVMAPI - struct
type GetAllSecurityTagsAttachedToVMAPI struct {
	*api.BaseAPI
}

// NewGetAllAttachedToVM - Generates a new GetAllSecurityTagsAttachedToVMAPI object.
func NewGetAllAttachedToVM(vmID string) *GetAllSecurityTagsAttachedToVMAPI {
	this := new(GetAllSecurityTagsAttachedToVMAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodGet, "/api/2.0/services/securitytags/vm/"+vmID, nil, new(SecurityTags))
	return this
}

// GetResponse - returns the ResponseObject from GetAllSecurityTagsAttachedToVMAPI
func (getAPI GetAllSecurityTagsAttachedToVMAPI) GetResponse() *SecurityTags {
	return getAPI.ResponseObject().(*SecurityTags)
}
