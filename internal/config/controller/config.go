package controller

// Config of controller contains data needed to run the controller app.
type Config struct {
	CSV    string `koanf:"csv"`    // testcase file path
	Client string `koanf:"client"` // target client name
}
