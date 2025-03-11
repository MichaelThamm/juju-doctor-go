package juju

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// GetJujuStatusOutput fetches the Juju status for the specified model.
func GetJujuStatusOutput(model string) (map[string]any, error) {

	args := []string{"status", "--format=json"}
	if model != "" {
		args = append(args, "--model", model)
	}
	cmd := exec.Command("juju", args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error executing [juju %s]: %w", strings.Join(args, " "), err)
	}

	var data map[string]any
	if err := json.Unmarshal(output, &data); err != nil {
		return nil, fmt.Errorf("error parsing Juju status JSON: %w", err)
	}

	return data, nil
}
