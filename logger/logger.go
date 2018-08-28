package logger

import (
	"fmt"
	"io"
	"os"

	loggerMw "github.com/dictyBase/go-middlewares/middlewares/logrus"

	"github.com/johntdyer/slackrus"
	"github.com/sirupsen/logrus"
	cli "gopkg.in/urfave/cli.v1"
)

// GetLoggerMiddleware gets a net/http compatible instance of logrus
func GetLoggerMiddleware(c *cli.Context) (*loggerMw.Logger, error) {
	var logger *loggerMw.Logger
	var w io.Writer
	if c.IsSet("log-file") {
		w, err := os.Open(c.String("log-file"))
		if err != nil {
			return logger, fmt.Errorf("could not open log file for writing %s", err)
		}
		w = io.MultiWriter(fw, os.Stderr)
	} else {
		w = os.Stderr
	}
	if c.String("log-format") == "json" {
		logger = loggerMw.NewJSONFileLogger(w)
	} else {
		logger = loggerMw.NewFileLogger(w)
	}
	return logger, nil
}

// GetAppLogger gets a configured instance of logrus for logging by application
func GetAppLogger(c *cli.Context) (*logrus.Logger, error) {
	log := logrus.New()
	if c.IsSet("app-log") {
		w, err := os.Create(c.String("app-log"))
		if err != nil {
			return log, fmt.Errorf("error in opening log file %s", err)
		}
		log.Out = w
	} else {
		log.Out = os.Stderr
	}

	switch c.String("app-log-fmt") {
	case "text":
		log.Formatter = &logrus.TextFormatter{
			TimestampFormat: "02/Jan/2006:15:04:05",
		}
	case "json":
		log.Formatter = &logrus.JSONFormatter{
			TimestampFormat: "02/Jan/2006:15:04:05",
		}
	}
	l := c.String("app-log-level")
	switch l {
	case "debug":
		log.Level = logrus.DebugLevel
	case "info":
		log.Level = logrus.InfoLevel
	case "warn":
		log.Level = logrus.WarnLevel
	case "error":
		log.Level = logrus.ErrorLevel
	case "fatal":
		log.Level = logrus.FatalLevel
	case "panic":
		log.Level = logrus.PanicLevel
	}
	// Set up hook
	lh := make(logrus.LevelHooks)
	for _, h := range c.StringSlice("hooks") {
		switch h {
		case "slack":
			lh.Add(&slackrus.SlackrusHook{
				HookURL:        c.String("slack-url"),
				AcceptedLevels: slackrus.LevelThreshold(log.Level),
				IconEmoji:      ":skull:",
			})
		default:
			continue
		}
	}
	log.Hooks = lh
	return log, nil
}
