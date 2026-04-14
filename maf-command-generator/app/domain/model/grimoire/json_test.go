package grimoire

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestGrimoireValidateAllWithRealJSONTestcases(t *testing.T) {
	caseDir := "testcases"
	dirEntries, err := os.ReadDir(caseDir)
	if err != nil {
		t.Fatalf("read testcases dir: %v", err)
	}

	var files []string
	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		files = append(files, entry.Name())
	}
	sort.Strings(files)
	if len(files) == 0 {
		t.Fatalf("no json testcase found in %s", caseDir)
	}

	for _, file := range files {
		file := file
		t.Run(file, func(t *testing.T) {
			path := filepath.Join(caseDir, file)
			entity := NewGrimoireEntity(path)
			if err := entity.Load(); err != nil {
				t.Fatalf("load testcase %s: %v", file, err)
			}

			allErrs := entity.ValidateAll(testDBMaster{})
			switch {
			case strings.HasSuffix(file, ".ok.json"):
				if len(allErrs) != 0 {
					t.Fatalf("expected no validation errors for %s, got %#v", file, allErrs)
				}
			case strings.HasSuffix(file, ".ng.json"):
				if len(allErrs) == 0 {
					t.Fatalf("expected validation errors for %s, got none", file)
				}
			default:
				t.Fatalf("testcase file must end with .ok.json or .ng.json: %s", file)
			}
		})
	}
}
