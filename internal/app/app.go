package app

import (
	"flag"
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"io"
	"log/slog"
	"os"
	"strings"
)

type logger interface {
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
}

type App struct {
	defaultConfigPath  string
	overrideConfigPath string
	mergedConfigPath   string

	log logger
}

func NewApp() *App {
	a := &App{}
	a.init()

	return a
}

func (a *App) init() {
	flag.StringVar(&a.overrideConfigPath, "override", "", "override config path")
	flag.StringVar(&a.mergedConfigPath, "merged", "", "merged config path")
	flag.StringVar(&a.defaultConfigPath, "default", "", "default config path")
	logFormat := flag.String("log-format", "standard", "log format")
	flag.Parse()

	var logHandler slog.Handler
	switch *logFormat {
	case "json":
		logHandler = slog.NewJSONHandler(os.Stdout, nil)
	default:
		logHandler = slog.NewTextHandler(os.Stdout, nil)
	}

	a.log = slog.New(logHandler)
	if a.defaultConfigPath == "" || a.overrideConfigPath == "" || a.mergedConfigPath == "" {
		a.log.Error("config path, override config path and merged config path flag values are required")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func (a *App) Run() {
	def, err := os.Open(a.defaultConfigPath)
	if err != nil {
		a.log.Error("could not open default file", "error", err.Error(), "path", a.defaultConfigPath)
		os.Exit(1)
	}

	overr, err := os.Open(a.overrideConfigPath)
	if err != nil {
		a.log.Error("could not open override file", "error", err.Error(), "path", a.defaultConfigPath)
		os.Exit(1)
	}

	splitConf := strings.Split(a.defaultConfigPath, ".")
	confExt := splitConf[len(splitConf)-1]
	var mergeBytes []byte

	switch confExt {
	case "toml":
		mergeBytes, err = a.processToml(def, overr)
		if err != nil {
			a.log.Error("could not process toml", "error", err.Error())
			os.Exit(1)
		}
	default:
		a.log.Error("file extension not supported", "extension", confExt)
		os.Exit(0)
	}

	if err = os.WriteFile(a.mergedConfigPath, mergeBytes, 0644); err != nil {
		a.log.Error("could not write merged file:", err)
		os.Exit(1)
	}

	a.log.Info("merged config created successfully", "path", a.mergedConfigPath)
}

func (a *App) processToml(def io.Reader, overr io.Reader) ([]byte, error) {
	var defToml map[string]interface{}
	var overrToml map[string]interface{}

	if err := toml.NewDecoder(def).Decode(&defToml); err != nil {
		return nil, fmt.Errorf("coulod not decode default toml: %w", err)
	}

	if err := toml.NewDecoder(overr).Decode(&overrToml); err != nil {
		return nil, fmt.Errorf("coulod not decode override toml: %w", err)
	}

	for oKey, oVal := range overrToml {
		oValMap, ok := oVal.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("could not convert override toml field %s to map", oKey)
		}

		if dVal, ok := defToml[oKey]; ok {
			dValMap, ok := dVal.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("could not convert default toml field %s to map", oKey)
			}

			for overrideKey, overrideVal := range oValMap {
				if _, ok := dValMap[overrideKey]; ok {
					a.log.Info("overriding default config",
						"default", fmt.Sprintf("%s=%v", overrideKey, dValMap[overrideKey]),
						"override", fmt.Sprintf("%s=%v", overrideKey, overrideVal),
					)
					dValMap[overrideKey] = overrideVal
				}
			}

			defToml[oKey] = dValMap
		}
	}

	tomlOut, err := toml.Marshal(defToml)
	if err != nil {
		return nil, fmt.Errorf("could not marshal default toml: %w", err)
	}

	return tomlOut, nil

}
