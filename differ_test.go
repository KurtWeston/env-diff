package main

import (
	"testing"
)

func TestCompareEnvFiles(t *testing.T) {
	t.Run("detect missing variables", func(t *testing.T) {
		env1 := &EnvFile{
			Path: "file1",
			Vars: map[string]EnvVar{
				"KEY1": {Key: "KEY1", Value: "val1", Line: 1},
			},
		}
		env2 := &EnvFile{
			Path: "file2",
			Vars: map[string]EnvVar{
				"KEY1": {Key: "KEY1", Value: "val1", Line: 1},
				"KEY2": {Key: "KEY2", Value: "val2", Line: 2},
			},
		}

		diff := CompareEnvFiles(env1, env2, "file1", "file2")

		if len(diff.Missing) != 1 {
			t.Errorf("expected 1 missing var, got %d", len(diff.Missing))
		}
		if diff.Missing[0].Key != "KEY2" {
			t.Errorf("expected missing KEY2, got %s", diff.Missing[0].Key)
		}
	})

	t.Run("detect extra variables", func(t *testing.T) {
		env1 := &EnvFile{
			Path: "file1",
			Vars: map[string]EnvVar{
				"KEY1": {Key: "KEY1", Value: "val1", Line: 1},
				"KEY2": {Key: "KEY2", Value: "val2", Line: 2},
			},
		}
		env2 := &EnvFile{
			Path: "file2",
			Vars: map[string]EnvVar{
				"KEY1": {Key: "KEY1", Value: "val1", Line: 1},
			},
		}

		diff := CompareEnvFiles(env1, env2, "file1", "file2")

		if len(diff.Extra) != 1 {
			t.Errorf("expected 1 extra var, got %d", len(diff.Extra))
		}
	})

	t.Run("detect mismatched values", func(t *testing.T) {
		env1 := &EnvFile{
			Path: "file1",
			Vars: map[string]EnvVar{
				"KEY1": {Key: "KEY1", Value: "val1", Line: 1},
			},
		}
		env2 := &EnvFile{
			Path: "file2",
			Vars: map[string]EnvVar{
				"KEY1": {Key: "KEY1", Value: "val2", Line: 1},
			},
		}

		diff := CompareEnvFiles(env1, env2, "file1", "file2")

		if len(diff.Mismatched) != 1 {
			t.Errorf("expected 1 mismatch, got %d", len(diff.Mismatched))
		}
	})

	t.Run("identical files", func(t *testing.T) {
		env1 := &EnvFile{
			Path: "file1",
			Vars: map[string]EnvVar{
				"KEY1": {Key: "KEY1", Value: "val1", Line: 1},
			},
		}
		env2 := &EnvFile{
			Path: "file2",
			Vars: map[string]EnvVar{
				"KEY1": {Key: "KEY1", Value: "val1", Line: 1},
			},
		}

		diff := CompareEnvFiles(env1, env2, "file1", "file2")

		if diff.HasDifferences() {
			t.Error("expected no differences for identical files")
		}
		if len(diff.Matching) != 1 {
			t.Errorf("expected 1 matching var, got %d", len(diff.Matching))
		}
	})
}

func TestHasDifferences(t *testing.T) {
	tests := []struct {
		name string
		diff *DiffResult
		want bool
	}{
		{"no diffs", &DiffResult{}, false},
		{"has missing", &DiffResult{Missing: []EnvVar{{Key: "K"}}}, true},
		{"has extra", &DiffResult{Extra: []EnvVar{{Key: "K"}}}, true},
		{"has mismatch", &DiffResult{Mismatched: []MismatchedVar{{Key: "K"}}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.diff.HasDifferences(); got != tt.want {
				t.Errorf("HasDifferences() = %v, want %v", got, tt.want)
			}
		})
	}
}