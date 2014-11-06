package setting

import (
	"fmt"
	"github.com/Unknwon/com"
	"github.com/Unknwon/goconfig"
	"github.com/weisd/goapi/modules/log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

var (
	Cfg                *goconfig.ConfigFile
	Domain             string
	HttpAddr, HttpPort string
	CustomPath         string
	ConfRootPath       string

	// Log settings.
	LogRootPath string
	LogModes    []string
	LogConfigs  []string
)

var logLevels = map[string]string{
	"Trace":    "0",
	"Debug":    "1",
	"Info":     "2",
	"Warn":     "3",
	"Error":    "4",
	"Critical": "5",
}

func ExecPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	p, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	return p, nil
}

// WorkDir returns absolute path of work directory.
func WorkDir() (string, error) {
	execPath, err := ExecPath()
	return path.Dir(strings.Replace(execPath, "\\", "/", -1)), err
}

// NewConfigContext initializes configuration context.
// NOTE: do not print any log except error.
func NewConfigContext() {
	workDir, err := WorkDir()
	if err != nil {
		log.Fatal(4, "Fail to get work directory: %v", err)
	}
	ConfRootPath = path.Join(workDir, "conf")

	Cfg, err = goconfig.LoadConfigFile(path.Join(workDir, "conf/app.ini"))
	if err != nil {
		log.Fatal(4, "Fail to parse 'conf/app.ini': %v", err)
	}

	CustomPath = os.Getenv("GOGS_CUSTOM")
	if len(CustomPath) == 0 {
		CustomPath = path.Join(workDir, "custom")
	}

	cfgPath := path.Join(CustomPath, "conf/app.ini")
	if com.IsFile(cfgPath) {
		if err = Cfg.AppendFiles(cfgPath); err != nil {
			log.Fatal(4, "Fail to load custom 'conf/app.ini': %v", err)
		}
	} else {
		log.Warn("No custom 'conf/app.ini' found, please go to '/install'")
	}
	Domain = Cfg.MustValue("server", "DOMAIN", "localhost")
	HttpAddr = Cfg.MustValue("server", "HTTP_ADDR", "0.0.0.0")
	HttpPort = Cfg.MustValue("server", "HTTP_PORT", "3000")

}

func newLogService() {
	// log.Info("%s %s", AppName, AppVer)

	// Get and check log mode.
	LogModes = strings.Split(Cfg.MustValue("log", "MODE", "console"), ",")
	fmt.Println("%v", LogModes)
	LogConfigs = make([]string, len(LogModes))
	for i, mode := range LogModes {
		mode = strings.TrimSpace(mode)
		modeSec := "log." + mode
		fmt.Println(modeSec, i)
		if _, err := Cfg.GetSection(modeSec); err != nil {
			log.Fatal(4, "Unknown log mode: %s", mode)
		}

		// Log level.
		levelName := Cfg.MustValueRange("log."+mode, "LEVEL", "Trace",
			[]string{"Trace", "Debug", "Info", "Warn", "Error", "Critical"})
		level, ok := logLevels[levelName]
		if !ok {
			log.Fatal(4, "Unknown log level: %s", levelName)
		}

		// Generate log configuration.
		switch mode {
		case "console":
			LogConfigs[i] = fmt.Sprintf(`{"level":%s}`, level)
		case "file":
			logPath := Cfg.MustValue(modeSec, "FILE_NAME", path.Join(LogRootPath, "gogs.log"))
			os.MkdirAll(path.Dir(logPath), os.ModePerm)
			LogConfigs[i] = fmt.Sprintf(
				`{"level":%s,"filename":"%s","rotate":%v,"maxlines":%d,"maxsize":%d,"daily":%v,"maxdays":%d}`, level,
				logPath,
				Cfg.MustBool(modeSec, "LOG_ROTATE", true),
				Cfg.MustInt(modeSec, "MAX_LINES", 1000000),
				1<<uint(Cfg.MustInt(modeSec, "MAX_SIZE_SHIFT", 28)),
				Cfg.MustBool(modeSec, "DAILY_ROTATE", true),
				Cfg.MustInt(modeSec, "MAX_DAYS", 7))
		case "conn":
			LogConfigs[i] = fmt.Sprintf(`{"level":%s,"reconnectOnMsg":%v,"reconnect":%v,"net":"%s","addr":"%s"}`, level,
				Cfg.MustBool(modeSec, "RECONNECT_ON_MSG"),
				Cfg.MustBool(modeSec, "RECONNECT"),
				Cfg.MustValueRange(modeSec, "PROTOCOL", "tcp", []string{"tcp", "unix", "udp"}),
				Cfg.MustValue(modeSec, "ADDR", ":7020"))
		case "smtp":
			LogConfigs[i] = fmt.Sprintf(`{"level":%s,"username":"%s","password":"%s","host":"%s","sendTos":"%s","subject":"%s"}`, level,
				Cfg.MustValue(modeSec, "USER", "example@example.com"),
				Cfg.MustValue(modeSec, "PASSWD", "******"),
				Cfg.MustValue(modeSec, "HOST", "127.0.0.1:25"),
				Cfg.MustValue(modeSec, "RECEIVERS", "[]"),
				Cfg.MustValue(modeSec, "SUBJECT", "Diagnostic message from serve"))
		case "database":
			LogConfigs[i] = fmt.Sprintf(`{"level":%s,"driver":"%s","conn":"%s"}`, level,
				Cfg.MustValue(modeSec, "DRIVER"),
				Cfg.MustValue(modeSec, "CONN"))
		}

		log.NewLogger(Cfg.MustInt64("log", "BUFFER_LEN", 10000), mode, LogConfigs[i])
		log.Info("Log Mode: %s(%s)", strings.Title(mode), levelName)
	}
}

func NewServices() {
	newLogService()
}
