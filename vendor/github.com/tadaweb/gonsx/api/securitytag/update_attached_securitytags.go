package securitytag

import (
	"github.com/tadaweb/gonsx/api"
	"net/http"
)

// UpdateAttachedSecurityTagsAPI - struct
type UpdateAttachedSecurityTagsAPI struct {
	*api.BaseAPI
}

// NewUpdateAttachedTags - Generates a NewUpdateAttachedSecurityTagsAPI object.
func NewUpdateAttachedTags(vmID string, securityTagPayload *AttachmentList) *UpdateAttachedSecurityTagsAPI {
	this := new(UpdateAttachedSecurityTagsAPI)
	endpointURL := "/api/2.0/services/securitytags/vm/" + vmID + "?action=ASSIGN_TAGS"
	this.BaseAPI = api.NewBaseAPI(http.MethodPost, endpointURL, securityTagPayload, nil)
	return this
}

// GetResponse returns the ResponseObject from CreateSecurityTagAPI
func (updateAPI UpdateAttachedSecurityTagsAPI) GetResponse() string {
	return updateAPI.ResponseObject().(string)
}
