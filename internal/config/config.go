// Package config provides unified configuration loading from YAML files
// with CLI override support using dot-notation keys.
//
// Example config.yaml:
//
//	accounts:
//	  database:
//	    url: sqlite:///accounts.db
//	  address:
//	    host: 127.0.0.1
//	    port: 8000
//
// CLI override: ./accounts --address.port=8080 --database.url=postgres://...
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	ct "phoenix/internal/common/types"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// Config holds service-scoped configuration loaded from YAML
// with CLI overrides applied.
type Config struct {
	k *koanf.Koanf
}

// Usage returns the help/usage string for the binary.
func Usage() string {
	return strings.TrimSpace(`
Usage: <binary> [options]

Options:
  -c, --config <path>   Path to configuration file (default: config.yaml)
                        Can also be set via CONFIG environment variable.
  -h, --help            Show this help message.

Config overrides (dot-notation):
  --section.key=value    Override any config value at runtime.
                        Example: --address.port=9090
`) + "\n"
}

// Load reads the YAML config file, extracts the section matching
// serviceName, applies CLI overrides from --key=value arguments,
// and returns a scoped Config.
//
// The config file path is resolved in this order:
//  1. -c/--config CLI flag
//  2. CONFIG environment variable
//  3. Default: "config.yaml"
//
// If -h or --help is passed, Load prints usage and calls os.Exit(0).
func Load(serviceName string) (*Config, error) {
	configPath, helpRequested, overrides := parseArgs(os.Args[1:])
	if helpRequested {
		fmt.Print(Usage())
		os.Exit(0)
	}

	// Resolve config path: CLI flag > CONFIG env > default
	if configPath == "" || configPath == "config.yaml" {
		if envPath := os.Getenv("CONFIG"); envPath != "" {
			configPath = envPath
		}
	}
	if configPath == "" {
		configPath = "config.yaml"
	}

	// Load YAML file
	k := koanf.New(".")
	if err := k.Load(file.Provider(configPath), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("load config %q: %w\nhint: use -c/--config flag or CONFIG env var to specify config path", configPath, err)
	}

	// Scope to service section
	svcK := koanf.New(".")
	if raw := k.Get(serviceName); raw != nil {
		if m, ok := raw.(map[string]any); ok {
			svcK.Load(confmap.Provider(m, ""), nil)
		}
	}

	// Apply CLI overrides (keys are relative to service section,
	// values already type-inferred by parseArgs).
	for key, val := range overrides {
		svcK.Set(key, val)
	}

	return &Config{k: svcK}, nil
}

// Unmarshal deserializes a config subtree into a struct.
// Use "" for the root of the service config, or a dotted path
// for a sub-section (e.g., "database").
func (c *Config) Unmarshal(path string, out any) error {
	return c.k.Unmarshal(path, out)
}

// String returns a config value as string (dot-notation key).
func (c *Config) String(key string) string {
	return c.k.String(key)
}

// Int returns a config value as int (dot-notation key).
func (c *Config) Int(key string) int {
	return c.k.Int(key)
}

// Exists checks whether a key is present in the config.
func (c *Config) Exists(key string) bool {
	return c.k.Exists(key)
}

// Currencies parses the "currensies" section into a slice of Currency.
//
// Expected YAML shape:
//
//	currensies:
//	  - ETH:                    # native coin → route = itself
//	  - USD(T):
//	      - network: ETH
//	        token: USDT
func (c *Config) Currencies() ([]ct.Currency, error) {
	raw := c.k.Get("currensies")
	if raw == nil {
		return nil, nil
	}

	list, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("currensies: expected a list, got %T", raw)
	}

	var currencies []ct.Currency
	for _, item := range list {
		m, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("currensies: expected map, got %T", item)
		}
		for symbol, val := range m {
			curr := ct.Currency{Symbol: symbol}
			if val == nil {
				// Native coin: its only route is itself.
				curr.Routes = []ct.CurrencyRoute{
					{Network: symbol, Token: symbol},
				}
			} else {
				routesList, ok := val.([]interface{})
				if !ok {
					return nil, fmt.Errorf("currensies.%s: expected list of routes, got %T", symbol, val)
				}
				for _, r := range routesList {
					rm, ok := r.(map[string]interface{})
					if !ok {
						return nil, fmt.Errorf("currensies.%s: expected route map, got %T", symbol, r)
					}
					curr.Routes = append(curr.Routes, ct.CurrencyRoute{
						Network: rm["network"].(string),
						Token:   rm["token"].(string),
					})
				}
			}
			currencies = append(currencies, curr)
		}
	}
	return currencies, nil
}

// ---------------------------------------------------------------------------
// Internal helpers

// parseArgs extracts --config/-c path, detects --help/-h, and collects
// --key=value overrides from raw args.
func parseArgs(args []string) (configPath string, helpRequested bool, overrides map[string]any) {
	configPath = "" // empty means "not set via CLI"
	overrides = make(map[string]any)

	for i := 0; i < len(args); i++ {
		arg := args[i]

		// Handle --help and -h
		if arg == "--help" || arg == "-h" {
			helpRequested = true
			continue
		}

		// Handle --config=path and -c=path (equals form)
		if strings.HasPrefix(arg, "--config=") {
			configPath = strings.TrimPrefix(arg, "--config=")
			continue
		}
		if strings.HasPrefix(arg, "-c=") {
			configPath = strings.TrimPrefix(arg, "-c=")
			continue
		}

		// Handle --config path and -c path (space-separated form)
		if (arg == "--config" || arg == "-c") && i+1 < len(args) {
			configPath = args[i+1]
			i++ // skip next arg
			continue
		}

		// Handle --key=value overrides
		if !strings.HasPrefix(arg, "--") {
			continue
		}
		arg = strings.TrimPrefix(arg, "--")

		key, val, found := strings.Cut(arg, "=")
		if !found {
			continue
		}

		if key == "config" {
			configPath = val
		} else {
			overrides[key] = inferType(val)
		}
	}

	return configPath, helpRequested, overrides
}

// inferType tries to parse a string value into the most appropriate Go type:
// bool > int > float > string. This ensures that CLI overrides
// like --port=8080 produce integer values, not strings.
func inferType(s string) any {
	// bool
	if b, err := strconv.ParseBool(s); err == nil {
		return b
	}
	// int64
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i
	}
	// float64
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return s
}
