package logs

// Reset turns the values back to initial state.
func (l *Logs) Reset() {
	l.datastore = make(map[int]interface{}, 0)
	l.logs = make(map[int]interface{}, 0)
}
