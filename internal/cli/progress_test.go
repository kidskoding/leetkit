package cli

import (
	"strings"
	"testing"

	"github.com/kidskoding/leetkit/internal/progress"
)

func TestProgress(t *testing.T) {
	env := tempEnv(t)
	code, out, _ := run(t, env, "progress")
	if code != exitOK || !strings.Contains(out, "Overall") || !strings.Contains(out, "0 solved · 0 test runs") {
		t.Fatalf("exit %d\n%s", code, out)
	}
	if strings.Contains(out, "last activity") {
		t.Error("no activity yet, but a last-activity time was shown")
	}

	store, err := progress.Open(env.Paths.ProgressFile())
	if err != nil {
		t.Fatal(err)
	}
	for id, passed := range map[string]bool{"enumerate-01": true, "zip-01": false} {
		if _, err := store.Record(id, passed); err != nil {
			t.Fatal(err)
		}
	}
	code, out, _ = run(t, env, "progress")
	for _, want := range []string{"Iteration and Unpacking", "1/9", "11%", "1 solved · 2 test runs · last activity"} {
		if code != exitOK || !strings.Contains(out, want) {
			t.Errorf("missing %q (exit %d):\n%s", want, code, out)
		}
	}
}
