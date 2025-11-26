package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseEnvFile(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("parse valid env file", func(t *testing.T) {
		filePath := filepath.Join(tmpDir, "test.env")
		content := "DB_HOST=localhost\nDB_PORT=5432\n# Comment line\nAPI_KEY=secret123\n"
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		envFile, err := ParseEnvFile(filePath)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(envFile.Vars) != 3 {
			t.Errorf("expected 3 vars, got %d", len(envFile.Vars))
		}

		if envFile.Vars["DB_HOST"].Value != "localhost" {
			t.Errorf("expected DB_HOST=localhost, got %s", envFile.Vars["DB_HOST"].Value)
		}
	})

	t.Run("file not found", func(t *testing.T) {
		_, err := ParseEnvFile("nonexistent.env")
		if err == nil {
			t.Error("expected error for nonexistent file")
		}
	})

	t.Run("parse with inline comments", func(t *testing.T) {
		filePath := filepath.Join(tmpDir, "comments.env")
		content := "API_KEY=secret # production key\nDEBUG=true\n"
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		envFile, err := ParseEnvFile(filePath)
		if err != nil {
			t.Fatal(err)
		}

		if envFile.Vars["API_KEY"].Comment != "production key" {
			t.Errorf("expected comment 'production key', got '%s'", envFile.Vars["API_KEY"].Comment)
		}
	})

	t.Run("parse quoted values", func(t *testing.T) {
		filePath := filepath.Join(tmpDir, "quoted.env")
		content := "MSG=\"hello world\"\nPATH='usr/bin'\n"
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		envFile, err := ParseEnvFile(filePath)
		if err != nil {
			t.Fatal(err)
		}

		if envFile.Vars["MSG"].Value != "hello world" {
			t.Errorf("expected unquoted value, got %s", envFile.Vars["MSG"].Value)
		}
	})
}

func TestParseLine(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		wantKey string
		wantVal string
		wantCmt string
	}{
		{"simple", "KEY=value", "KEY", "value", ""},
		{"with comment", "KEY=value # comment", "KEY", "value", "comment"},
		{"quoted", "KEY=\"val\"", "KEY", "val", ""},
		{"no equals", "INVALID", "", "", ""},
		{"empty value", "KEY=", "KEY", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, val, cmt := parseLine(tt.line)
			if key != tt.wantKey || val != tt.wantVal || cmt != tt.wantCmt {
				t.Errorf("parseLine(%q) = (%q, %q, %q), want (%q, %q, %q)",
					tt.line, key, val, cmt, tt.wantKey, tt.wantVal, tt.wantCmt)
			}
		})
	}
}

func TestUnquote(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"\"quoted\"", "quoted"},
		{"'quoted'", "quoted"},
		{"unquoted", "unquoted"},
		{"\"", "\""},
	}

	for _, tt := range tests {
		got := unquote(tt.input)
		if got != tt.want {
			t.Errorf("unquote(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}