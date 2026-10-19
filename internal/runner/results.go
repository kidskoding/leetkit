package runner

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

type Status int

const (
	StatusPassed Status = iota
	StatusFailed
	// StatusError means tests could not run, usually a syntax or import error in the learner's code.
	StatusError
	StatusTimedOut
	StatusCanceled
)

func (s Status) String() string {
	switch s {
	case StatusPassed:
		return "passed"
	case StatusFailed:
		return "failed"
	case StatusError:
		return "error"
	case StatusTimedOut:
		return "timed out"
	case StatusCanceled:
		return "canceled"
	}
	return fmt.Sprintf("Status(%d)", int(s))
}

type Outcome string

const (
	OutcomePassed  Outcome = "passed"
	OutcomeFailed  Outcome = "failed"
	OutcomeError   Outcome = "error"
	OutcomeSkipped Outcome = "skipped"
)

type Case struct {
	Name    string
	Outcome Outcome
	// Message is the short reason, such as the assertion or the exception line.
	Message string
	// Details is pytest's traceback for the case.
	Details string
}

type Result struct {
	Status   Status
	Cases    []Case
	Duration time.Duration
	// Output is pytest's console output, kept so learners can inspect full tracebacks.
	Output string
	// Problem explains StatusError, StatusTimedOut, and StatusCanceled results.
	Problem string
}

func (r *Result) Count(o Outcome) int {
	n := 0
	for _, c := range r.Cases {
		if c.Outcome == o {
			n++
		}
	}
	return n
}

type junitDetail struct {
	Message string `xml:"message,attr"`
	Text    string `xml:",chardata"`
}

type junitCase struct {
	Name    string       `xml:"name,attr"`
	Failure *junitDetail `xml:"failure"`
	Error   *junitDetail `xml:"error"`
	Skipped *junitDetail `xml:"skipped"`
}

type junitSuite struct {
	Cases []junitCase `xml:"testcase"`
}

// parseJUnit accepts both a <testsuites> root (pytest >= 5.1) and a bare <testsuite>.
func parseJUnit(data []byte) ([]Case, error) {
	var root struct {
		Suites []junitSuite `xml:"testsuite"`
		Cases  []junitCase  `xml:"testcase"`
	}
	if err := xml.Unmarshal(data, &root); err != nil {
		return nil, fmt.Errorf("parse pytest report: %w", err)
	}
	raw := root.Cases
	for _, s := range root.Suites {
		raw = append(raw, s.Cases...)
	}
	cases := make([]Case, 0, len(raw))
	for _, jc := range raw {
		c := Case{Name: jc.Name, Outcome: OutcomePassed}
		switch {
		case jc.Error != nil:
			c.Outcome, c.Details = OutcomeError, jc.Error.Text
			c.Message = errorSummary(jc.Error)
		case jc.Failure != nil:
			c.Outcome, c.Details = OutcomeFailed, jc.Failure.Text
			c.Message = jc.Failure.Message
		case jc.Skipped != nil:
			c.Outcome, c.Message = OutcomeSkipped, jc.Skipped.Message
		}
		cases = append(cases, c)
	}
	return cases, nil
}

// errorSummary prefers pytest's "E " lines, which carry the actual exception,
// over generic messages such as "collection failure".
func errorSummary(d *junitDetail) string {
	var lines []string
	for line := range strings.SplitSeq(d.Text, "\n") {
		if rest, ok := strings.CutPrefix(line, "E "); ok {
			lines = append(lines, strings.TrimPrefix(rest, "  "))
		}
	}
	if len(lines) > 0 {
		return strings.Join(lines, "\n")
	}
	return d.Message
}
