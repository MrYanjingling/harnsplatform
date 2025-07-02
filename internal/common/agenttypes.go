package common

type AgentType int64

const (
	AgentTypeNone AgentType = iota
	AgentTypeModbus
)

var AgentTypeToString = map[AgentType]string{
	AgentTypeModbus: "modbus",
}

var StringToAgentType = map[string]AgentType{
	"modbus": AgentTypeModbus,
}

// func (dt AccessMode) MarshalJSON() ([]byte, error) {
// 	if s, ok := ReadWritePropertyToString[dt]; ok {
// 		return json.Marshal(s)
// 	}
// 	return nil, fmt.Errorf("unknown accessMode %d", dt)
// }
//
// func (dt *AccessMode) UnmarshalJSON(bytes []byte) error {
// 	var s string
// 	if err := json.Unmarshal(bytes, &s); err != nil {
// 		return err
// 	}
//
// 	v, ok := StringToReadWriteProperty[s]
// 	if !ok {
// 		return fmt.Errorf("unknown accessMode %s", s)
// 	}
// 	*dt = v
// 	return nil
// }
