package utils

var (
	ServerName = "AZALF"
	ServerDesc = "Adam's Zillenial Arch Linux Flavor Server"
	ServerPort = ":9999"
	ServerVers = "0.2.9"

	DevelopmentFile = "C:\\Users\\adam\\.config\\azalf\\.azalf.yml"
	Debug           = false

	// config is going to be used throughout the server to serve out configurations
	// This should be constant and only changed if ./azalf update is called or at
	// at the start of the server.
	AzalfConfig *Config

	// Debug values to describe the Debug type and message
	ERROR   = "error"
	WARN    = "warning"
	SUCCESS = "success"
)
