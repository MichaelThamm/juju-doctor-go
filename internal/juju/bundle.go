package juju

import (
	"gopkg.in/yaml.v3"
	"fmt"
	"os/exec"
	"strings"
)

// GetJujuBundleOutput fetches the Juju bundle for the specified model.
func GetJujuBundleOutput(model string) (any, error) {

	args := []string{"export-bundle"}
	if model != "" {
		args = append(args, "--model", model)
	}
	cmd := exec.Command("juju", args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("Error executing [juju %s]: %w", strings.Join(args, " "), err)
	}

	var singleDoc map[string]any
	var multiDoc []map[string]any
	if err := yaml.Unmarshal(output, &singleDoc); err == nil {
		return singleDoc, nil // It's a single YAML document
	} else if err := yaml.Unmarshal(output, &multiDoc); err == nil {
		return multiDoc, nil // It's a multi-document YAML
	} else {
		return nil, fmt.Errorf("error parsing Juju bundle YAML: %w", err)
	}
}
