package logs

import (
	"github.com/f24-cse535/pbft/pkg/models"
)

// GetLog is returns a log by its index.
func (l *Logs) GetLog(index int) *models.Log {
	if value, ok := l.logs[index]; ok {
		return value
	}

	return nil
}
