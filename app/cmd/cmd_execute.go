package cmd

import (
	"fmt"
	"strings"

	"github.com/spyzhov/healthy/step"
)

func (app Application) execute(groups []string, steps []string) int {
	if len(groups) == 0 {
		groups = app.StepGroups.Names()
	}
	totals := map[step.Status]int{
		step.Success: 0,
		step.Warning: 0,
		step.Error:   0,
	}
	for _, gName := range groups {
		printed := false
		group := app.StepGroups.Get(gName)
		names := steps
		if len(names) == 0 {
			names = group.Names()
		}
		for _, name := range names {
			current := group.Get(name)
			if current == nil {
				continue
			}
			if !app.Config.Quiet {
				if !printed {
					fmt.Printf("%s:\n", gName)
					printed = true
				}

				fmt.Printf("  %s -> ...", name)
			}
			res := current.Call()
			if !app.Config.Quiet {
				if app.Config.Verbose {
					fmt.Printf("\r  %s -> %s!\n%s\n", name, res.Status, shift(res.Message, 5))
				} else {
					fmt.Printf("\r  %s -> %s!\n", name, res.Status)
				}
			}
			totals[res.Status]++
		}
	}
	if !app.Config.Quiet {
		fmt.Printf(`
Total
    success: %d
    warinig: %d
     errors: %d
`, totals[step.Success], totals[step.Warning], totals[step.Error])
	}
	if totals[step.Error] > 0 {
		return 1
	}
	return 0
}

func shift(str string, size int) string {
	small := strings.Repeat(" ", size-3)
	delim := strings.Repeat(" ", size)
	parts := strings.Split(str, "\n")
	return small + "-> " + strings.Join(parts, "\n"+delim)
}
