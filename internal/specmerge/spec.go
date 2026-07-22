// Package specmerge loads authored eval spec files and normalizes them into a
// single assertion list for downstream execution.
package specmerge

import (
	"encoding/json"
	"fmt"
	"os"
)

// Assertion is one structured check an eval runs.
type Assertion struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

// Eval is one authored eval: free-form expectations plus structured
// assertions.
type Eval struct {
	Name         string      `json:"name"`
	Expectations []string    `json:"expectations"`
	Assertions   []Assertion `json:"assertions"`
}

// SpecFile is a decoded spec file.
type SpecFile struct {
	Evals []Eval `json:"evals"`
}

// Load reads, decodes, and normalizes a spec file.
func Load(path string) (*SpecFile, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read spec: %w", err)
	}
	var f SpecFile
	if err := json.Unmarshal(raw, &f); err != nil {
		return nil, fmt.Errorf("decode spec %s: %w", path, err)
	}
	f.normalize()
	return &f, nil
}

// normalize expands each eval's authored expectations into leading llm
// assertions so downstream execution handles one list.
func (f *SpecFile) normalize() {
	for i := range f.Evals {
		e := &f.Evals[i]
		merged := make([]Assertion, 0, len(e.Expectations)+len(e.Assertions))
		for _, want := range e.Expectations {
			merged = append(merged, Assertion{Kind: "llm", Value: want})
		}
		merged = append(merged, e.Assertions...)
		e.Assertions = merged
		e.Expectations = nil
	}
}
