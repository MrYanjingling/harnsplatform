package v1

import (
	"harnsplatform/internal/biz"
)

type Brokers struct {
	Name            string                              `json:"name,omitempty"`
	ThingTypeId     *string                             `json:"parentTypeId,omitempty"`
	Description     string                              `json:"description,omitempty"`
	Characteristics map[string]*biz.Characteristics     `json:"characteristics,omitempty"`
	PropertySets    map[string]map[string]*biz.Property `json:"propertySets,omitempty"`
	Combination     []string                            `json:"combination,omitempty"`
	*biz.Meta       `json:",inline"`
}
