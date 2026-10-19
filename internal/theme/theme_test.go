package theme

import (
	"reflect"
	"strings"
	"testing"
)

func TestPalettesAreComplete(t *testing.T) {
	for name, p := range palettes {
		v := reflect.ValueOf(p)
		for i := range v.NumField() {
			if v.Field(i).IsNil() {
				t.Errorf("%s: %s is not set", name, v.Type().Field(i).Name)
			}
		}
	}
}

func TestLoad(t *testing.T) {
	th, err := Load("")
	if err != nil || th.Name != Default {
		t.Fatalf("Load(\"\") = %q, %v; want default", th.Name, err)
	}
	for _, name := range Names() {
		if _, err := Load(name); err != nil {
			t.Errorf("Load(%q): %v", name, err)
		}
	}
	_, err = Load("solarized")
	if err == nil || !strings.Contains(err.Error(), "available: gruvbox") {
		t.Errorf("Load(unknown) err = %v; want list of themes", err)
	}
}
