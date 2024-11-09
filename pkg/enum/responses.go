package enum

// responses are the text values that are returned to the clients.
const (
	RespBadBalance       = "declined the client does not have enough balance"
	RespSuccess          = "transaction submitted"
	RespNotEnoughServers = "not enough servers are available"
)
