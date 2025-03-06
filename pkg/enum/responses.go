package enum

// responses are the text values that are returned to the clients.
const (
	RespNotEnoughBalance   = "declined the client does not have enough balance"
	RespSuccessTransaction = "transaction submitted"
	RespNotEnoughServers   = "not enough servers are available"
	RespSystemFailed       = "system is unavailable"
)
