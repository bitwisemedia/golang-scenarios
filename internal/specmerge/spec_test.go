package specmerge

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadNormalizesExpectations(t *testing.T) {
	spec := `{
	  "evals": [
	    {
	      "name": "greeting",
	      "expectations": ["is polite", "mentions the user by name"],
	      "assertions": [{"kind": "contains", "value": "hello"}]
	    }
	  ]
	}`
	path := filepath.Join(t.TempDir(), "spec.json")
	if err := os.WriteFile(path, []byte(spec), 0o600); err != nil {
		t.Fatal(err)
	}

	f, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(f.Evals) != 1 {
		t.Fatalf("got %d evals, want 1", len(f.Evals))
	}
	e := f.Evals[0]
	if len(e.Assertions) != 3 {
		t.Fatalf("got %d assertions, want 3", len(e.Assertions))
	}
	if e.Assertions[0].Kind != "llm" || e.Assertions[0].Value != "is polite" {
		t.Errorf("first assertion = %+v, want expanded llm expectation", e.Assertions[0])
	}
	if e.Assertions[2].Kind != "contains" {
		t.Errorf("last assertion = %+v, want the authored assertion", e.Assertions[2])
	}
	if e.Expectations != nil {
		t.Errorf("expectations not cleared after normalize: %v", e.Expectations)
	}
}
