package main

import (
	"fmt"
	"sort"

	"github.com/fatih/color"
)

type DiffResult struct {
	File1        string
	File2        string
	Missing      []EnvVar
	Extra        []EnvVar
	Mismatched   []MismatchedVar
	Matching     []EnvVar
}

type MismatchedVar struct {
	Key   string
	Var1  EnvVar
	Var2  EnvVar
}

func CompareEnvFiles(env1, env2 *EnvFile, name1, name2 string) *DiffResult {
	result := &DiffResult{
		File1:      name1,
		File2:      name2,
		Missing:    []EnvVar{},
		Extra:      []EnvVar{},
		Mismatched: []MismatchedVar{},
		Matching:   []EnvVar{},
	}

	for key, var2 := range env2.Vars {
		if var1, exists := env1.Vars[key]; exists {
			if var1.Value != var2.Value {
				result.Mismatched = append(result.Mismatched, MismatchedVar{Key: key, Var1: var1, Var2: var2})
			} else {
				result.Matching = append(result.Matching, var1)
			}
		} else {
			result.Missing = append(result.Missing, var2)
		}
	}

	for key, var1 := range env1.Vars {
		if _, exists := env2.Vars[key]; !exists {
			result.Extra = append(result.Extra, var1)
		}
	}

	sort.Slice(result.Missing, func(i, j int) bool { return result.Missing[i].Key < result.Missing[j].Key })
	sort.Slice(result.Extra, func(i, j int) bool { return result.Extra[i].Key < result.Extra[j].Key })
	sort.Slice(result.Mismatched, func(i, j int) bool { return result.Mismatched[i].Key < result.Mismatched[j].Key })

	return result
}

func (d *DiffResult) HasDifferences() bool {
	return len(d.Missing) > 0 || len(d.Extra) > 0 || len(d.Mismatched) > 0
}

func (d *DiffResult) Print(noColor, quiet bool) {
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	if noColor {
		color.NoColor = true
	}

	if !quiet {
		if len(d.Missing) > 0 {
			fmt.Printf("\n%s\n", red("Missing variables (in "+d.File2+" but not in "+d.File1+"):"))
			for _, v := range d.Missing {
				fmt.Printf("  %s %s (line %d)\n", red("-"), cyan(v.Key), v.Line)
				if v.Comment != "" {
					fmt.Printf("    # %s\n", v.Comment)
				}
			}
		}

		if len(d.Extra) > 0 {
			fmt.Printf("\n%s\n", yellow("Extra variables (in "+d.File1+" but not in "+d.File2+"):"))
			for _, v := range d.Extra {
				fmt.Printf("  %s %s (line %d)\n", yellow("+"), cyan(v.Key), v.Line)
			}
		}

		if len(d.Mismatched) > 0 {
			fmt.Printf("\n%s\n", yellow("Mismatched values:"))
			for _, m := range d.Mismatched {
				fmt.Printf("  %s\n", cyan(m.Key))
				fmt.Printf("    %s: %s (line %d)\n", d.File1, m.Var1.Value, m.Var1.Line)
				fmt.Printf("    %s: %s (line %d)\n", d.File2, m.Var2.Value, m.Var2.Line)
			}
		}
	}

	fmt.Printf("\n%s\n", green("Summary:"))
	fmt.Printf("  Total in %s: %d\n", d.File1, len(d.Extra)+len(d.Mismatched)+len(d.Matching))
	fmt.Printf("  Total in %s: %d\n", d.File2, len(d.Missing)+len(d.Mismatched)+len(d.Matching))
	fmt.Printf("  %s: %d\n", red("Missing"), len(d.Missing))
	fmt.Printf("  %s: %d\n", yellow("Extra"), len(d.Extra))
	fmt.Printf("  %s: %d\n", yellow("Mismatched"), len(d.Mismatched))
	fmt.Printf("  %s: %d\n", green("Matching"), len(d.Matching))

	if d.HasDifferences() {
		fmt.Printf("\n%s\n", red("❌ Differences found!"))
	} else {
		fmt.Printf("\n%s\n", green("✓ Files are in sync!"))
	}
}
