package logs

// Logs is a memory type that stores the node's logs and datastore.
type Logs struct {
	logs      map[int]interface{}
	datastore map[int]interface{}
}

// NewLogs returns a new logs instance.
func NewLogs() *Logs {
	return &Logs{
		logs:      make(map[int]interface{}),
		datastore: make(map[int]interface{}),
	}
}
