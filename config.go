package godrone

var DefaultConfig = Config{
	LogLevel:      "debug",
	LogTimeFormat: "15:04:05.999999",
	MotorboardTTY: "/dev/ttyO0",
	NavboardTTY:   "/dev/ttyO1",
	HttpAPIPort:   80,
}

// @TODO embedd controller config
type Config struct {
	LogLevel      string
	LogTimeFormat string
	MotorboardTTY string
	NavboardTTY   string
	HttpAPIPort   int
}
