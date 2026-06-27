package system

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v3"

	"github.com/iftsoft/driver/config"
)

type AppParams struct {
	AppName string
	DevName string
	LogPath string
	CfgPath string
}

func RunBootstrap(cfg *config.AppConfig) (*slog.Logger, error) {
	// get application params
	params := GetAppParams()
	cfgFile := filepath.Join(params.CfgPath, params.DevName+".yaml")
	logFile := filepath.Join(params.LogPath, params.DevName+".log")

	// read application config
	err := ReadYamlFile(cfgFile, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read app config %s: %w", cfgFile, err)
	}

	// prepare logs folder
	if params.LogPath != "" {
		err := os.MkdirAll(params.LogPath, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("failed to create log directory %s: %w", params.LogPath, err)
		}
	}

	// create log file rotator
	var logWriter io.Writer
	rotator := GetRotator(cfg.Logger, logFile)
	if rotator != nil {
		logWriter = rotator
	} else {
		logWriter = os.Stdout
	}

	// prepare logger options
	level := slog.LevelDebug
	_ = level.UnmarshalText([]byte(cfg.Logger.Level))
	opts := &slog.HandlerOptions{
		Level: level,
	}
	// create logger
	logger := slog.New(slog.NewTextHandler(logWriter, opts))
	slog.SetDefault(logger)
	logger.Info("Run application bootstrap", slog.Any("params", params), slog.Any("logger", cfg.Logger))

	return logger, nil
}

func GetAppParams() AppParams {
	params := AppParams{
		AppName: filepath.Base(os.Args[0]),
	}
	flag.StringVar(&params.DevName, "dev_name", "", "unique device name")
	flag.StringVar(&params.LogPath, "log_path", "", "path to logs folder")
	flag.StringVar(&params.CfgPath, "cfg_path", "", "path to configs folder")
	flag.Parse()
	if params.DevName == "" { // call without params
		params.DevName = params.AppName
	}
	if params.LogPath == "" {
		params.LogPath = fmt.Sprintf(".%clog", os.PathSeparator)
	}
	if params.CfgPath == "" {
		params.CfgPath = fmt.Sprintf(".%ccfg", os.PathSeparator)
	}

	return params
}

func GetRotator(cfg config.LoggerConfig, fileName string) *lumberjack.Logger {
	if fileName == "" {
		return nil
	}
	if cfg.MaxSize == 0 {
		cfg.MaxSize = 10
	}
	if cfg.MaxBackups == 0 {
		cfg.MaxBackups = 3
	}
	if cfg.MaxAge == 0 {
		cfg.MaxAge = 30
	}
	rotator := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    cfg.MaxSize,    // Megabytes before rotating
		MaxBackups: cfg.MaxBackups, // Max number of old log files to retain
		MaxAge:     cfg.MaxAge,     // Days to retain old log files
		Compress:   cfg.Compress,   // Compress rotated files using gzip
	}
	return rotator
}

func ReadYamlFile(name string, cfg interface{}) error {
	data, err := os.ReadFile(name)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return err
	}
	return nil
}
