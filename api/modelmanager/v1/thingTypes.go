package v1

import (
	"harnsplatform/internal/biz"
)

type ThingTypes struct {
	Name            string                              `json:"name,omitempty"`
	ParentTypeId    string                              `json:"parentTypeId,omitempty"`
	Description     string                              `json:"description,omitempty"`
	Characteristics map[string]*biz.Characteristics     `json:"characteristics,omitempty"`
	PropertySets    map[string]map[string]*biz.Property `json:"propertySets,omitempty"`
	*biz.Meta       `json:",inline"`
}
