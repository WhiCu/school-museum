package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/kdl"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var (
	ConfigPath = "PATH_CONFIG"

	k = koanf.New(".")
)

type filetype int

const (
	YAML filetype = iota
	JSON
	TOML
	KDL
)

func FileType(s string) (filetype, error) {
	switch s {
	case "yaml", "yml":
		return YAML, nil
	case "json":
		return JSON, nil
	case "toml":
		return TOML, nil
	case "kdl":
		return KDL, nil
	default:
		return 0, fmt.Errorf("unsupported file type: %s", s)
	}
}

func parser(ft filetype) koanf.Parser {
	switch ft {
	case YAML:
		return yaml.Parser()
	case JSON:
		return json.Parser()
	case TOML:
		return toml.Parser()
	case KDL:
		return kdl.Parser()
	default:
		panic("unsupported file type")
	}
}

func Load[T any](ft filetype) (*T, error) {
	path := os.Getenv(ConfigPath)

	if err := check(path); err != nil {
		return nil, err
	}

	if err := k.Load(file.Provider(path), parser(ft)); err != nil {
		return nil, err
	}

	if err := k.Load(env.Provider(".", env.Opt{
		TransformFunc: func(k, v string) (string, any) {
			var newKey string

			switch {
			case strings.HasPrefix(k, "SERVER_"):
				newKey = strings.Replace(strings.ToLower(k), "server_", "server.", 1)
			case strings.HasPrefix(k, "DB_"):
				newKey = strings.Replace(strings.ToLower(k), "db_", "storage.", 1)
			case strings.HasPrefix(k, "LOG_"):
				newKey = strings.Replace(strings.ToLower(k), "log_", "logger.", 1)
			default:
				return "", nil
			}

			if strings.Contains(v, " ") {
				return newKey, strings.Split(v, " ")
			}

			return newKey, v
		},
	}), nil); err != nil {
		return nil, err
	}

	var cfg T
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func check(path string) error {
	if path == "" {
		return fmt.Errorf("%s is not set", ConfigPath)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("%s does not exist", path)
	}

	return nil
}

func MustLoad[T any](ft filetype) *T {
	cfg, err := Load[T](ft)
	if err != nil {
		panic(err)
	}
	return cfg
}
