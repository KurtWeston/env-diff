package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type EnvVar struct {
	Key     string
	Value   string
	Line    int
	Comment string
}

type EnvFile struct {
	Path string
	Vars map[string]EnvVar
}

func ParseEnvFile(path string) (*EnvFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	envFile := &EnvFile{
		Path: path,
		Vars: make(map[string]EnvVar),
	}

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		key, value, comment := parseLine(line)
		if key != "" {
			envFile.Vars[key] = EnvVar{
				Key:     key,
				Value:   value,
				Line:    lineNum,
				Comment: comment,
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return envFile, nil
}

func parseLine(line string) (key, value, comment string) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", ""
	}

	key = strings.TrimSpace(parts[0])
	valPart := parts[1]

	commentIdx := strings.Index(valPart, "#")
	if commentIdx != -1 && !isInsideQuotes(valPart, commentIdx) {
		comment = strings.TrimSpace(valPart[commentIdx+1:])
		valPart = valPart[:commentIdx]
	}

	value = strings.TrimSpace(valPart)
	value = unquote(value)

	return key, value, comment
}

func isInsideQuotes(s string, pos int) bool {
	inQuote := false
	for i := 0; i < pos && i < len(s); i++ {
		if s[i] == '"' || s[i] == '\'' {
			inQuote = !inQuote
		}
	}
	return inQuote
}

func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
