package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/canonical/starform/starform"
	"github.com/canonical/starlark/starlark"
)

type WriterLogger struct {
	Writer       io.Writer
	MinimumLevel starform.LogLevel
}

func (l *WriterLogger) Log(ctx context.Context, entry starform.LogEntry) {
	if entry.Level < l.MinimumLevel {
		return
	}
	// Format the log message properly
	logMsg := fmt.Sprintf("[%v] %s %s", entry.EventName, entry.Path, entry.Message)

	fmt.Fprintln(l.Writer, logMsg)
}

func loadScriplets(ctx context.Context) (*starform.ScriptSet, error) {
	cache, err := starform.NewDefaultCache(&starform.DefaultCacheOptions{
		MaxSize: 100, // The maximum number of cache entries before automatic eviction starts.
	})
	if err != nil {
		return nil, err
	}
	scriptset, err := starform.NewScriptSet(&starform.ScriptSetOptions{
		App: &starform.AppObject{
			Name: "jujudoctor",
		},
		Cache: cache,
		Logger: &WriterLogger{
			Writer:       os.Stdout,
			MinimumLevel: starform.PrintLevel,
		},
		// RequiredSafety: starlark.CPUSafe | starlark.IOSafe | starlark.MemSafe,
		// MaxAllocs:      1024,
		// MaxSteps:       1000,
	})
	if err != nil {
		return nil, err
	}

	fs := os.DirFS(".")
	sources, err := starform.LoadDirSources(ctx, &starform.LoadDirSourcesOptions{
		FS:   fs,
		Root: "scriptlets",
	})
	if err != nil {
		return nil, err
	}

	// Load the listed scripts from disk and prepare for event handling.
	err = scriptset.LoadSources(ctx, sources)
	if err != nil {
		return nil, err
	}

	return scriptset, nil

}

func toStarlarkDict(rawData map[string]interface{}) (starlark.Value, error) {
	starlarkDict := starlark.NewDict(len(rawData))
	for key, value := range rawData {
		starlarkValue, err := toStarlarkValue(value)
		if err != nil {
			return nil, fmt.Errorf("error converting key %q: %v", key, err)
		}
		starlarkDict.SetKey(starlark.String(key), starlarkValue)
	}

	return starlarkDict, nil
}

func toStarlarkValue(value interface{}) (starlark.Value, error) {
	switch v := value.(type) {
	case string:
		return starlark.String(v), nil
	case int:
		return starlark.MakeInt(v), nil
	case int64:
		return starlark.MakeInt64(v), nil
	case float64:
		return starlark.Float(v), nil
	case bool:
		return starlark.Bool(v), nil
	case map[string]interface{}:
		starlarkDict := starlark.NewDict(len(v))
		for key, val := range v {
			starlarkVal, err := toStarlarkValue(val)
			if err != nil {
				return nil, fmt.Errorf("error converting key %q: %v", key, err)
			}
			if err := starlarkDict.SetKey(starlark.String(key), starlarkVal); err != nil {
				return nil, fmt.Errorf("error setting key %q: %v", key, err)
			}
		}
		return starlarkDict, nil
	case []interface{}:
		starlarkList := make([]starlark.Value, len(v))
		for i, elem := range v {
			starlarkElem, err := toStarlarkValue(elem)
			if err != nil {
				return nil, fmt.Errorf("error converting list element %d: %v", i, err)
			}
			starlarkList[i] = starlarkElem
		}
		return starlark.NewList(starlarkList), nil
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("unsupported type: %T", value)
	}
}

func getJujuStatusOutput() (starlark.Value, error) {
	cmd := exec.Command("juju", "status", "--model", "test", "--format=json")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(output, &data); err != nil {
		return nil, err
	}

	return toStarlarkDict(data)
}

func main() {
	log.SetFlags(0)
	ctx := context.Background()

	// load available scriplets from "./scriplets" directory
	scriptSet, err := loadScriplets(ctx)
	if err != nil {
		log.Fatalf("error loading scriplets %s", err.Error())
	}

	// Fire pre-defined events
	jujuStatusObj, err := getJujuStatusOutput()
	if err != nil {
		log.Fatalf("error getting juju status outputr %s", err.Error())
	}

	// Fire an event that is expected to succeed
	err = scriptSet.Handle(ctx, &starform.EventObject{
		Name:  "status_ready",
		Attrs: starlark.StringDict{"input": jujuStatusObj, "error": starlark.False},
	})
	if err != nil {
		log.Fatalf("failed to handle event: %v", err)
	}

	// Fire an event that is expected to fail
	err = scriptSet.Handle(ctx, &starform.EventObject{
		Name:  "status_ready",
		Attrs: starlark.StringDict{"input": jujuStatusObj, "error": starlark.True},
	})
	if err != nil {
		log.Fatalf("[status_ready] ERROR: %v", err)
	}

}
