package godrone

import (
	//"github.com/felixge/godrone/apis/http"
	"github.com/felixge/godrone/controller/drivers"
	"github.com/felixge/godrone/controller"
	"github.com/felixge/godrone/log"
	"os"
	"time"
)

func NewFirmware(c Config) (*Firmware, error) {
	log, err := log.NewLogger(c.LogLevel, c.LogTimeFormat, os.Stdout)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	lap := start
	log.Info("Initializing firmware")

	log.Debug("Initializing navboard at TTY: %s", c.NavboardTTY)
	navboard, err := drivers.NewNavboard(c.NavboardTTY, log)
	if err != nil {
		return nil, log.Emergency("Could not initialize navboard: %s", err)
	}
	log.Debug("Initialized navboard, took: %s", time.Since(lap))

	lap = time.Now()
	log.Debug("Initializing motorboard at TTY: %s", c.MotorboardTTY)
	motorboard, err := drivers.NewMotorboard(c.MotorboardTTY, log)
	if err != nil {
		return nil, log.Emergency("Could not initialize motorboard: %s", err)
	}
	log.Debug("Initialized navboard, took: %s", time.Since(lap))

	//lap = time.Now()
	//log.Debug("Initializing http api on port: %d", c.HttpAPIPort)
	//httpApi, err := http.NewHttpAPI(c.HttpAPIPort, , log)
	//if err != nil {
		//return nil, log.Emergency("Could not initialize http api: %s", err)
	//}
	//log.Debug("Initialized http api, took: %s", time.Since(lap))


	firmware := &Firmware{
		config:     &c,
		log:        log,
		//navboard:   navboard,
		//motorboard: motorboard,
		//httpApi:    httpApi,
	}
	log.Info("Initialized firmware, took: %s", time.Since(start))
	return firmware, nil
}

type Firmware struct {
	config     *Config
	log        log.Logger
	//navboard   *drivers.Navboard
	//motorboard *drivers.Motorboard
	//httpApi    *http.HttpAPI
}

func (f *Firmware) Run() error {
	f.log.Info("Starting main loop")
	f.log.Debug("Serving http api")
	return f.httpApi.Serve()
}
