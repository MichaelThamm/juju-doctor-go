package utils

import (
	"context"
	"fmt"
	"io"

	"github.com/canonical/starform/starform"
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
