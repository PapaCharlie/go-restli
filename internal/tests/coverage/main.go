package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

type TestEvent struct {
	Action  string
	Package string
	Test    string
	Output  string
}

const suite = "github.com/PapaCharlie/go-restli/v2/internal/tests/suite"
const rootTest = "TestGoRestli"

func main() {
	skippedTests := make(map[string]map[string]bool)
	deliberatelySkippedTests := make(map[string]int)

	decoder := json.NewDecoder(os.Stdin)
	for decoder.More() {
		var e TestEvent
		err := decoder.Decode(&e)
		if err != nil {
			panic(err)
		}

		var resource, operation string
		if e.Package == suite && strings.HasPrefix(e.Test, rootTest) {
			resource, operation, _ = strings.Cut(strings.TrimPrefix(e.Test, rootTest+"/"), "/")
			if operation != "" {
				if skippedTests[resource] == nil {
					skippedTests[resource] = make(map[string]bool)
				}
			}
		}

		switch e.Action {
		case "output":
			if operation != "" && strings.Contains(e.Output, "GORESTLI_SKIPPED") {
				deliberatelySkippedTests[resource]++
				skippedTests[resource][operation] = true
			}
			fmt.Print(e.Output)
		case "skip", "pass", "fail":
			if operation == "" || !strings.Contains(operation, "/") {
				continue
			}
			if _, ok := skippedTests[resource][operation]; ok {
				continue
			}
			skippedTests[resource][operation] = e.Action == "skip"
		}
	}

	var longestResource int
	resources := make([]string, 0, len(skippedTests))
	for r := range skippedTests {
		resources = append(resources, r)
		if l := len(r); l > longestResource {
			longestResource = l
		}
	}
	sort.Strings(resources)

	fmt.Print("\n\n")
	fmt.Println("REST.LI TEST SUITE COVERAGE")
	for _, r := range resources {
		operations := skippedTests[r]
		var skipped float64
		for _, s := range operations {
			if s {
				skipped++
			}
		}
		skipped -= float64(deliberatelySkippedTests[r])
		fmt.Printf("%s:%s %3.0f%%\n", r, strings.Repeat(" ", longestResource-len(r)), 100*(1-skipped/float64(len(operations))))
	}
	fmt.Print("\n\n")
}
