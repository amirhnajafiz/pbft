package models

// TestSet is a row in the test-case file.
type TestSet struct {
	Index            string
	LiveServers      []string
	ByzantineServers []string
	Transactions     []*Transaction
}
