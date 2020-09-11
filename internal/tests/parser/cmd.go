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

const suite = "github.com/PapaCharlie/go-restli/internal/tests/suite"
const rootTest = "TestGoRestli"

func main() {
	skippedTests := make(map[string]map[string]bool)

	decoder := json.NewDecoder(os.Stdin)
	for decoder.More() {
		var e TestEvent
		err := decoder.Decode(&e)
		if err != nil {
			panic(err)
		}

		switch e.Action {
		case "output":
			fmt.Print(e.Output)
		case "skip", "pass", "fail":
			if e.Package != suite || !strings.HasPrefix(e.Test, rootTest) {
				continue
			}
			nameSegments := strings.Split(e.Test, "/")
			if len(nameSegments) < 3 {
				continue
			}
			resource, operation := nameSegments[1], nameSegments[2]
			if skippedTests[resource] == nil {
				skippedTests[resource] = make(map[string]bool)
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
		fmt.Printf("%s:%s %3.0f%%\n", r, strings.Repeat(" ", longestResource-len(r)), 100*(1-skipped/float64(len(operations))))
	}
	fmt.Print("\n\n")
}
