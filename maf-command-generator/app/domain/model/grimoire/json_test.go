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

	var dirs []string
	for _, entry := range dirEntries {
		if !entry.IsDir() {
			continue
		}
		dirs = append(dirs, entry.Name())
	}
	sort.Strings(dirs)
	if len(dirs) == 0 {
		t.Fatalf("no testcase dir found in %s", caseDir)
	}

	for _, dir := range dirs {
		t.Run(dir, func(t *testing.T) {
			path := filepath.Join(caseDir, dir)
			entity := NewGrimoireEntity(path)
			if err := entity.Load(); err != nil {
				t.Fatalf("load testcase %s: %v", dir, err)
			}

			allErrs := entity.ValidateAll(testDBMaster{})
			switch {
			case strings.HasSuffix(dir, ".ok"):
				if len(allErrs) != 0 {
					t.Fatalf("expected no validation errors for %s, got %#v", dir, allErrs)
				}
			case strings.HasSuffix(dir, ".ng"):
				if len(allErrs) == 0 {
					t.Fatalf("expected validation errors for %s, got none", dir)
				}
			default:
				t.Fatalf("testcase dir must end with .ok or .ng: %s", dir)
			}
		})
	}
}
