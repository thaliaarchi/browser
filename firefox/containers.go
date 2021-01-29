package firefox

// Containers contains containers contained in containers.json.
type Containers struct {
	Version           int64               `json:"version"` // i.e. 4
	LastUserContextID int64               `json:"lastUserContextId"`
	Identities        []ContainerIdentity `json:"identities"`
}

// ContainerIdentity is a container definition.
type ContainerIdentity struct {
	UserContextID int64  `json:"userContextId"`
	Public        bool   `json:"public"`
	Icon          string `json:"icon"`  // i.e. "circle"
	Color         string `json:"color"` // i.e. "blue"
	L10nID        string `json:"l10nID,omitempty"`
	AccessKey     string `json:"accessKey,omitempty"`
	TelemetryID   int64  `json:"telemetryId,omitempty"`
	Name          string `json:"name,omitempty"`
}

// ParseContainers parses the containers.json file in a Firefox profile.
func ParseContainers(filename string) (*Containers, error) {
	var containers Containers
	if err := parseJSON(filename, &containers); err != nil {
		return nil, err
	}
	return &containers, nil
}
