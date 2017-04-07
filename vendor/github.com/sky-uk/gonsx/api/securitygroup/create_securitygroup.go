package securitygroup

import (
	"github.com/sky-uk/gonsx/api"
	"net/http"
)

// CreateSecurityGroupAPI api object
type CreateSecurityGroupAPI struct {
	*api.BaseAPI
}

// NewCreate returns a new object of CreateSecurityGroupAPI.
func NewCreate(scopeID, securityGroupName, setOperator, criteriaOperator, criteriaKey, criteriaValue, criteria string) *CreateSecurityGroupAPI {
	this := new(CreateSecurityGroupAPI)
	requestPayload := new(SecurityGroup)
	requestPayload.Name = securityGroupName

	dynamicCriteria := DynamicCriteria{
		Operator: criteriaOperator,
		Key:      criteriaKey,
		Value:    criteriaValue,
		Criteria: criteria,
	}
	dynamicCriteriaList := []DynamicCriteria{dynamicCriteria}

	dynamicSet := DynamicSet{
		Operator:        setOperator,
		DynamicCriteria: dynamicCriteriaList,
	}
	dynamicSetList := []DynamicSet{dynamicSet}

	requestPayload.DynamicMemberDefinition = &DynamicMemberDefinition{dynamicSetList}
	this.BaseAPI = api.NewBaseAPI(http.MethodPost, "/api/2.0/services/securitygroup/bulk/"+scopeID, requestPayload, new(string))
	return this
}

// GetResponse returns a ResponseObject of CreateSecurityGroupAPI.
func (ca CreateSecurityGroupAPI) GetResponse() string {
	return ca.ResponseObject().(string)
}
