package juju


import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	
	"github.com/tidwall/gjson"
)

// GetJujuShowUnitOutput fetches Juju show-unit and converts it to a Starlark object.
func GetJujuShowUnitOutput(model string) (map[string]any, error) {

	// Get all the unit names from the model
	jujuStatusObj, err := GetJujuStatusOutput(model)
	if err != nil {
		return nil, err
	}
	jsonData, err := json.Marshal(jujuStatusObj)
	if err != nil {
		return nil, fmt.Errorf("Error converting Juju status to JSON:", err)
	}
	unitNames := []string{}
	gjson.Get(string(jsonData), "applications").ForEach(func(appKey, appValue gjson.Result) bool {
		appValue.Get("units").ForEach(func(unitKey, _ gjson.Result) bool {
			unitNames = append(unitNames, unitKey.String())
			return true
		})
		return true
	})
	
	// Create a map of each unit's show-unit content
	results := make(map[string]any)
	for _, unitName := range unitNames {
		args := []string{"show-unit", unitName, "--format=json"}
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
			return nil, fmt.Errorf("error parsing Juju status JSON for unit %s: %w", unitName, err)
		}
		results[unitName] = data
	}

	return results, nil
}
