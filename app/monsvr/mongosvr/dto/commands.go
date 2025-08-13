package mongodto

type CommandDataDto struct {
	CommandCode int                    `json:"command_code"`
	Message     string                 `json:"message"`
	Parameters  map[string]interface{} `json:"parameters"`
}
