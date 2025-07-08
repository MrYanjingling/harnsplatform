package collector

import (
	pb "harnsplatform/api/modelmanager/v1"
	"sync"
)

type ModelManager struct {
	ac         pb.AgentsHTTPClient
	ttc        pb.ThingTypesHTTPClient
	tc         pb.ThingsHTTPClient
	things     *sync.Map
	thingTypes *sync.Map
	agents     *sync.Map
}

func NewModelManager(ac pb.AgentsHTTPClient, ttc pb.ThingTypesHTTPClient, tc pb.ThingsHTTPClient) *ModelManager {
	return &ModelManager{
		ac:         ac,
		ttc:        ttc,
		tc:         tc,
		things:     &sync.Map{},
		thingTypes: &sync.Map{},
		agents:     &sync.Map{},
	}
}

func (m *ModelManager) Init() {

}

func (m *ModelManager) GetAgents() *sync.Map {
	return m.agents
}
