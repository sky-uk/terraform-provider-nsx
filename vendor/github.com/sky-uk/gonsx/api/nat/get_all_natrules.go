package nat

import (
	"github.com/sky-uk/gonsx/api"
	"net/http"
)

// GetAllNatRulesAPI base object.
type GetAllNatRulesAPI struct {
	*api.BaseAPI
}

// NewGetAll returns the api object ofGetNatRuleAPI
func NewGetAll(edgeID string) *GetAllNatRulesAPI {
	this := new(GetAllNatRulesAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodGet, "/api/4.0/edges/"+edgeID+"/nat/config", nil, new(Nat))
	return this
}

// GetResponse returns ResponseObject ofGetNatRuleAPI
func (getAPI GetAllNatRulesAPI) GetResponse() Nat {
	return *getAPI.ResponseObject().(*Nat)
}
