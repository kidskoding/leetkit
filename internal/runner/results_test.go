package runner

import "testing"

func TestParseJUnitBareSuite(t *testing.T) {
	cases, err := parseJUnit([]byte(`<testsuite>
<testcase name="ok"/>
<testcase name="skip"><skipped message="not today"/></testcase>
<testcase name="bad"><failure message="assert 1 == 2">trace</failure></testcase>
</testsuite>`))
	if err != nil {
		t.Fatal(err)
	}
	want := []Case{
		{Name: "ok", Outcome: OutcomePassed},
		{Name: "skip", Outcome: OutcomeSkipped, Message: "not today"},
		{Name: "bad", Outcome: OutcomeFailed, Message: "assert 1 == 2", Details: "trace"},
	}
	if len(cases) != len(want) {
		t.Fatalf("got %+v", cases)
	}
	for i := range want {
		if cases[i] != want[i] {
			t.Errorf("case %d = %+v, want %+v", i, cases[i], want[i])
		}
	}
}

func TestParseJUnitErrorSummary(t *testing.T) {
	cases, err := parseJUnit([]byte(`<testsuites><testsuite>
<testcase name="test_exercise"><error message="collection failure">frame
E     File "starter.py", line 1
E   SyntaxError: expected ':'</error></testcase>
</testsuite></testsuites>`))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := cases[0].Message, "  File \"starter.py\", line 1\nSyntaxError: expected ':'"; got != want {
		t.Errorf("message = %q, want %q", got, want)
	}
}

func TestParseJUnitInvalid(t *testing.T) {
	if _, err := parseJUnit([]byte("<testsuite")); err == nil {
		t.Error("want error for malformed XML")
	}
}

func TestCleanMessage(t *testing.T) {
	in := "assert [1, 3] == [3, 1]\n  \n  At index 0 diff: 1 != 3\n  Use -v to get more diff\n"
	want := "assert [1, 3] == [3, 1]\n\n  At index 0 diff: 1 != 3"
	if got := cleanMessage(in); got != want {
		t.Errorf("cleanMessage = %q, want %q", got, want)
	}
}

func TestParseJUnitMarksCollectionErrors(t *testing.T) {
	cases, err := parseJUnit([]byte(`<testsuite>
<testcase classname="" name="test_exercise"><error message="collection failure">E   SyntaxError: x</error></testcase>
<testcase classname="test_exercise" name="test_a"><error message="fixture">E   boom</error></testcase>
</testsuite>`))
	if err != nil {
		t.Fatal(err)
	}
	if !cases[0].Collection || cases[1].Collection {
		t.Errorf("Collection flags = %v, %v; want true, false", cases[0].Collection, cases[1].Collection)
	}
}
