package juju

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/canonical/juju-doctor/internal/utils"

	starlark "github.com/canonical/starlark/starlark"
)

// GetJujuStatusOutput fetches Juju status and converts it to a Starlark object.
func GetJujuStatusOutput(model string) (starlark.Value, error) {

	args := []string{"status", "--format=json"}
	if model != "" {
		args = append(args, "--model", model)
	}
	cmd := exec.Command("juju", args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error executing [juju %s]: %w", strings.Join(args, " "), err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(output, &data); err != nil {
		return nil, fmt.Errorf("error parsing Juju status JSON: %w", err)
	}

	return utils.ToStarlarkDict(data)
}
