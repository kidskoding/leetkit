package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadMissingFileUsesDefaults(t *testing.T) {
	cfg, err := Load(filepath.Join(t.TempDir(), "absent.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if cfg != Default() {
		t.Errorf("got %+v, want defaults", cfg)
	}
}

func TestLoadValues(t *testing.T) {
	cfg, err := Load(writeConfig(t, `
python = "/opt/py/bin/python3"
editor = "code --wait"
timeout = "30s"
`))
	if err != nil {
		t.Fatal(err)
	}
	want := Config{Python: "/opt/py/bin/python3", Editor: "code --wait", Timeout: 30 * time.Second}
	if cfg != want {
		t.Errorf("got %+v, want %+v", cfg, want)
	}
}

func TestLoadErrors(t *testing.T) {
	tests := map[string]string{
		`python = `:         "parse",
		`pyhton = "x"`:      `unknown key "pyhton"`,
		`timeout = "-1s"`:   "timeout must be positive",
		`timeout = "never"`: "parse",
		`timeout = 5`:       "duration string",
	}
	for content, want := range tests {
		_, err := Load(writeConfig(t, content))
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%q: error %v, want substring %q", content, err, want)
		}
	}
}

func TestDefaultPathsHonorsXDG(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("XDG layout applies to Linux")
	}
	data := t.TempDir()
	t.Setenv("XDG_DATA_HOME", data)
	t.Setenv("XDG_CONFIG_HOME", "/cfg")
	p, err := DefaultPaths()
	if err != nil {
		t.Fatal(err)
	}
	if p.DataDir != filepath.Join(data, "leetkit") || p.ConfigFile != "/cfg/leetkit/config.toml" {
		t.Errorf("got %+v", p)
	}
	if p.VenvPython() != filepath.Join(data, "leetkit", "venv", "bin", "python") {
		t.Errorf("VenvPython = %s", p.VenvPython())
	}
}
