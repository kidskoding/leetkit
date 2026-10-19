// Package config loads user settings and resolves LeetKit's per-user directories.
package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/BurntSushi/toml"
)

const appName = "leetkit"

const DefaultTimeout = 10 * time.Second

type Config struct {
	// Python is an interpreter path or a command name looked up on PATH.
	Python string `toml:"python"`
	// Editor is a command line such as "code --wait"; it takes precedence over $VISUAL and $EDITOR.
	Editor  string        `toml:"editor"`
	Timeout time.Duration `toml:"timeout"`
}

func Default() Config {
	return Config{Timeout: DefaultTimeout}
}

// Load reads the config file at path. A missing file yields the defaults.
func Load(path string) (Config, error) {
	cfg := Default()
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return cfg, nil
	}
	if err != nil {
		return cfg, fmt.Errorf("read config: %w", err)
	}
	md, err := toml.Decode(string(data), &cfg)
	if err != nil {
		return cfg, fmt.Errorf("parse %s: %w", path, err)
	}
	if keys := md.Undecoded(); len(keys) > 0 {
		return cfg, fmt.Errorf("%s: unknown key %q", path, keys[0].String())
	}
	// TOML integers would silently decode as nanoseconds.
	if md.IsDefined("timeout") && md.Type("timeout") != "String" {
		return cfg, fmt.Errorf("%s: timeout must be a duration string, for example \"10s\"", path)
	}
	if cfg.Timeout <= 0 {
		return cfg, fmt.Errorf("%s: timeout must be positive, for example \"10s\"", path)
	}
	return cfg, nil
}

type Paths struct {
	ConfigFile string
	DataDir    string
}

func (p Paths) WorkspaceDir() string { return filepath.Join(p.DataDir, "workspace") }
func (p Paths) ProgressFile() string { return filepath.Join(p.DataDir, "progress.json") }
func (p Paths) VenvDir() string      { return filepath.Join(p.DataDir, "venv") }

// VenvPython is where a venv created at VenvDir keeps its interpreter.
func (p Paths) VenvPython() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(p.VenvDir(), "Scripts", "python.exe")
	}
	return filepath.Join(p.VenvDir(), "bin", "python")
}

func DefaultPaths() (Paths, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return Paths{}, fmt.Errorf("locate config directory: %w", err)
	}
	dataDir, err := userDataDir()
	if err != nil {
		return Paths{}, fmt.Errorf("locate data directory: %w", err)
	}
	return Paths{
		ConfigFile: filepath.Join(configDir, appName, "config.toml"),
		DataDir:    filepath.Join(dataDir, appName),
	}, nil
}

// userDataDir mirrors os.UserConfigDir for data, which the standard library lacks.
func userDataDir() (string, error) {
	switch runtime.GOOS {
	case "windows":
		if dir := os.Getenv("LocalAppData"); dir != "" {
			return dir, nil
		}
		return "", errors.New("%LocalAppData% is not set")
	case "darwin", "ios":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, "Library", "Application Support"), nil
	default:
		if dir := os.Getenv("XDG_DATA_HOME"); filepath.IsAbs(dir) {
			return dir, nil
		}
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, ".local", "share"), nil
	}
}
