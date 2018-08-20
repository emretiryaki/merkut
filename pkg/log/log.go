package log

import (
	"github.com/inconshreveable/log15"
	"gopkg.in/ini.v1"
				)

type Level int


const(
	LevelCrit Level=iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
)


var Root log15.Logger


type Logger interface {

	Debug(msg string, ctx ...interface{})
	Info(msg string, ctx ...interface{})
	Warn(msg string, ctx ...interface{})
	Error(msg string, ctx ...interface{})
	Crit(msg string, ctx ...interface{})

}
func init() {
	Root = log15.Root()
	Root.SetHandler(log15.DiscardHandler())
}

func New(logger string, ctx ...interface{}) Logger {
	params := append([]interface{}{"logger", logger}, ctx...)
	return Root.New(params...)
}

func Close() {

}

func ReadLoggingConfig(modes []string, logsPath string, cfg *ini.File) {
	Close()

	//defaultLevelName, _ := getLogLevelFromConfig("log", "info", cfg)
	//defaultFilters := getFilters(util.SplitString(cfg.Section("log").Key("filters").String()))
	//
	//handlers := make([]log15.Handler, 0)
	//
	//for _, mode := range modes {
	//	mode = strings.TrimSpace(mode)
	//	sec, err := cfg.GetSection("log." + mode)
	//	if err != nil {
	//		Root.Error("Unknown log mode", "mode", mode)
	//	}
	//
	//	// Log level.
	//	_, level := getLogLevelFromConfig("log."+mode, defaultLevelName, cfg)
	//	filters := getFilters(util.SplitString(sec.Key("filters").String()))
	//	format := getLogFormat(sec.Key("format").MustString(""))
	//
	//	var handler log15.Handler
	//
	//	// Generate log configuration.
	//	switch mode {
	//	case "console":
	//		handler = log15.StreamHandler(os.Stdout, format)
	//	case "file":
	//		fileName := sec.Key("file_name").MustString(filepath.Join(logsPath, "grafana.log"))
	//		os.MkdirAll(filepath.Dir(fileName), os.ModePerm)
	//		fileHandler := NewFileWriter()
	//		fileHandler.Filename = fileName
	//		fileHandler.Format = format
	//		fileHandler.Rotate = sec.Key("log_rotate").MustBool(true)
	//		fileHandler.Maxlines = sec.Key("max_lines").MustInt(1000000)
	//		fileHandler.Maxsize = 1 << uint(sec.Key("max_size_shift").MustInt(28))
	//		fileHandler.Daily = sec.Key("daily_rotate").MustBool(true)
	//		fileHandler.Maxdays = sec.Key("max_days").MustInt64(7)
	//		fileHandler.Init()
	//
	//		loggersToClose = append(loggersToClose, fileHandler)
	//		handler = fileHandler
	//	case "syslog":
	//		sysLogHandler := NewSyslog(sec, format)
	//
	//		loggersToClose = append(loggersToClose, sysLogHandler)
	//		handler = sysLogHandler
	//	}
	//
	//	for key, value := range defaultFilters {
	//		if _, exist := filters[key]; !exist {
	//			filters[key] = value
	//		}
	//	}
	//
	//	handler = LogFilterHandler(level, filters, handler)
	//	handlers = append(handlers, handler)
	//}

	//Root.SetHandler(log15.MultiHandler(handlers...))


}