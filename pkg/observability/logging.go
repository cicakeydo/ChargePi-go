package observability

import (
	"log/syslog"
	"path/filepath"

	"github.com/ChargePi/ChargePi-go/pkg/util"
	graylog "github.com/gemnasium/logrus-graylog-hook/v3"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"github.com/lorenzodonini/ocpp-go/ws"
	"github.com/orandin/lumberjackrus"
	log "github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
)

const (
	LogFileName = "chargepi.log"
	LogFileDir  = "/var/log/chargepi"
)

type LogType string

const (
	RemoteLogging  = LogType("remote")
	ConsoleLogging = LogType("console")
)

type LogFormat string

const (
	Syslog = LogFormat("syslog")
	Gelf   = LogFormat("gelf")
)

type Logging struct {
	LogTypes []Type `json:"logTypes,omitempty" yaml:"logTypes" mapstructure:"logTypes"`
}

type Type struct {
	Type    string  `json:"type,omitempty" yaml:"type" mapstructure:"type" validate:"required"` // remote, console
	Format  *string `json:"format,omitempty" yaml:"format" mapstructure:"format"`               // gelf, syslog, json, etc
	Address *string `json:"address,omitempty" yaml:"address" mapstructure:"address"`
}

// SetupLogging setup logs
func SetupLogging(logger *log.Logger, loggingConfig Logging, isDebug bool) {
	// Default logging settings
	logLevel := log.InfoLevel
	formatter := &log.JSONFormatter{}
	logger.SetFormatter(formatter)

	if isDebug {
		// Set underlying library loggers to debug level
		logLevel = log.DebugLevel
		ocppj.SetLogger(logger)
		ws.SetLogger(logger)
	}

	logger.SetLevel(logLevel)

	// Setup file logging
	fileLogging(logger, filepath.Join(LogFileDir, LogFileName))

	// Setup remote logging, if configured
	for _, logType := range loggingConfig.LogTypes {
		switch LogType(logType.Type) {
		case RemoteLogging:
			if util.IsNilInterfaceOrPointer(logType.Address) && util.IsNilInterfaceOrPointer(logType.Format) {
				remoteLogging(logger, *logType.Address, LogFormat(*logType.Format))
			}
		case ConsoleLogging:
		}
	}
}

func fileLogging(logger *log.Logger, fileName string) {
	hook, err := lumberjackrus.NewHook(
		&lumberjackrus.LogFile{
			Filename:   fileName,
			MaxSize:    200,
			MaxBackups: 20,
			MaxAge:     1,
			Compress:   false,
			LocalTime:  false,
		},
		logger.GetLevel(),
		logger.Formatter,
		nil,
	)

	if err != nil {
		panic(err)
	}

	logger.AddHook(hook)
}

// remoteLogging sends logs remotely to Graylog or any Syslog receiver.
func remoteLogging(logger *log.Logger, address string, format LogFormat) {
	var (
		hook log.Hook
		err  error
	)

	switch format {
	case Gelf:
		graylogHook := graylog.NewAsyncGraylogHook(address, map[string]interface{}{})
		hook = graylogHook
	case Syslog:
		hook, err = lSyslog.NewSyslogHook(
			"tcp",
			address,
			syslog.LOG_WARNING,
			"chargePi",
		)
	default:
		return
	}

	if err == nil {
		logger.AddHook(hook)
	}
}
