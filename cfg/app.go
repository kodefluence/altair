package cfg

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/kodefluence/altair/adapter"
	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/entity"
)

// defaultUpstreamTimeout caps how long a single upstream request may take.
// Mirrors module/router/usecase.defaultUpstreamTimeout — kept in sync so a
// missing app.yml field and a missing WithUpstreamTimeout option produce the
// same effective timeout. 30s is conservative enough for most APIs while
// still bounded enough to prevent goroutine leaks from hung upstreams.
const defaultUpstreamTimeout = 30 * time.Second

type app struct{}

// baseProxy carries the canonical (v1.0+) proxy block. Kept as its own type
// so an absent block (zero value) is distinguishable via the Host field.
type baseProxy struct {
	Host               string `yaml:"host"`
	UpstreamTimeout    string `yaml:"upstream_timeout"`
	MaxRequestBodySize string `yaml:"max_request_body_size"`
}

// parseByteSize accepts an integer (bytes) or an integer suffixed with B,
// KB, MB, or GB (binary multipliers — KB = 1024, MB = 1024², GB = 1024³).
// Empty string is "no cap" (returns 0). Anything else is a hard error so
// typos in the YAML can't silently disable the limit.
func parseByteSize(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}

	multiplier := int64(1)
	upper := strings.ToUpper(s)
	switch {
	case strings.HasSuffix(upper, "GB"):
		multiplier = 1024 * 1024 * 1024
		s = s[:len(s)-2]
	case strings.HasSuffix(upper, "MB"):
		multiplier = 1024 * 1024
		s = s[:len(s)-2]
	case strings.HasSuffix(upper, "KB"):
		multiplier = 1024
		s = s[:len(s)-2]
	case strings.HasSuffix(upper, "B"):
		s = s[:len(s)-1]
	}

	n, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid byte size %q: %w", s, err)
	}
	if n < 0 {
		return 0, fmt.Errorf("byte size must be non-negative, got %d", n)
	}
	return n * multiplier, nil
}

type baseAppConfig struct {
	Version string   `yaml:"version"`
	Plugins []string `yaml:"plugins"`
	Port    string   `yaml:"port"`
	// Legacy top-level proxy_host (pre-nested-block schema). Kept for
	// backward compatibility with existing deployments; new apps generated
	// by `altair new` use the `proxy:` block. If both are set, `proxy.host`
	// wins so operators can stage a migration without an outage.
	ProxyHost     string    `yaml:"proxy_host"`
	Proxy         baseProxy `yaml:"proxy"`
	AutoMigrate   bool      `yaml:"auto_migrate"`
	Authorization struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"authorization"`
}

func App() core.AppLoader {
	return &app{}
}

func (a *app) Compile(configPath string) (core.AppConfig, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	contents, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	compiledContents, err := compileTemplate(contents)
	if err != nil {
		return nil, err
	}

	var config baseAppConfig

	if err := yaml.Unmarshal(compiledContents, &config); err != nil {
		return nil, err
	}

	switch v := config.Version; v {
	case "1.0":
		var appConfigOption entity.AppConfigOption

		if config.Authorization.Username == "" {
			return nil, errors.New("config authorization `username` cannot be empty")
		}

		if config.Authorization.Password == "" {
			return nil, errors.New("config authorization `password` cannot be empty")
		}

		if config.Port == "" {
			appConfigOption.Port = 1304
		} else {
			port, err := strconv.Atoi(config.Port)
			if err != nil {
				return nil, err
			}

			appConfigOption.Port = port
		}

		// Resolve proxy host: nested block wins over the legacy top-level
		// field; default to www.local.host if neither is set.
		switch {
		case config.Proxy.Host != "":
			appConfigOption.ProxyHost = config.Proxy.Host
		case config.ProxyHost != "":
			appConfigOption.ProxyHost = config.ProxyHost
		default:
			appConfigOption.ProxyHost = "www.local.host"
		}

		if config.Proxy.UpstreamTimeout == "" {
			appConfigOption.UpstreamTimeout = defaultUpstreamTimeout
		} else {
			d, err := time.ParseDuration(config.Proxy.UpstreamTimeout)
			if err != nil {
				return nil, fmt.Errorf("invalid proxy.upstream_timeout %q: %w", config.Proxy.UpstreamTimeout, err)
			}
			appConfigOption.UpstreamTimeout = d
		}

		// max_request_body_size defaults to 0 (unlimited) so existing
		// deployments that don't set it continue to behave as before.
		bodyCap, err := parseByteSize(config.Proxy.MaxRequestBodySize)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy.max_request_body_size %q: %w", config.Proxy.MaxRequestBodySize, err)
		}
		appConfigOption.MaxRequestBodySize = bodyCap

		appConfigOption.Plugins = config.Plugins
		appConfigOption.AutoMigrate = config.AutoMigrate
		appConfigOption.Authorization.Username = config.Authorization.Username
		appConfigOption.Authorization.Password = config.Authorization.Password

		return adapter.AppConfig(entity.NewAppConfig(appConfigOption)), nil
	default:
		return nil, fmt.Errorf("undefined template version: %s for app.yaml", v)
	}
}
