package export

import (
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestValidateCleanPathAllowsMissingDirWithinOutputRoot(t *testing.T) {
	outputRoot := ioCleanupFixturePath(t, "safe_tree")
	cleanPath := filepath.Join(outputRoot, "data", "maf", "function", "generated", "missing")

	if err := validateCleanPath(outputRoot, cleanPath); err != nil {
		t.Fatalf("missing dir inside outputRoot should be allowed: %v", err)
	}
}

func TestValidateCleanPathAllowsJSONAndMcfunctionOnly(t *testing.T) {
	outputRoot := ioCleanupFixturePath(t, "safe_tree")
	cleanPath := filepath.Join(outputRoot, "data", "maf", "function", "generated")

	if err := validateCleanPath(outputRoot, cleanPath); err != nil {
		t.Fatalf("allowed extensions should pass validation: %v", err)
	}
}

func TestValidateCleanPathRejectsWhenOutputRootEmpty(t *testing.T) {
	cleanPath := filepath.Join(ioCleanupFixturePath(t, "safe_tree"), "data")

	err := validateCleanPath("", cleanPath)
	if err == nil {
		t.Fatal("expected error for empty outputRoot")
	}
	if !strings.Contains(err.Error(), "outputRoot must not be empty") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateCleanPathRejectsOutsideOutputRoot(t *testing.T) {
	outputRoot := ioCleanupFixturePath(t, "safe_tree")
	cleanPath := filepath.Join(ioCleanupFixturePath(t, "unsafe_tree"), "data", "maf", "function", "generated")

	err := validateCleanPath(outputRoot, cleanPath)
	if err == nil {
		t.Fatal("expected outside path to fail validation")
	}
	if !strings.Contains(err.Error(), "must be inside outputRoot") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateCleanPathRejectsPrefixLikeSibling(t *testing.T) {
	outputRoot := ioCleanupFixturePath(t, "safe_tree")
	cleanPath := outputRoot + "_side"

	err := validateCleanPath(outputRoot, cleanPath)
	if err == nil {
		t.Fatal("expected prefix-like sibling path to fail validation")
	}
	if !strings.Contains(err.Error(), "must be inside outputRoot") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateCleanPathRejectsTraversalOutsideOutputRoot(t *testing.T) {
	outputRoot := ioCleanupFixturePath(t, "safe_tree")
	cleanPath := filepath.Join(outputRoot, "..", "unsafe_tree", "data", "maf", "function", "generated")

	err := validateCleanPath(outputRoot, cleanPath)
	if err == nil {
		t.Fatal("expected traversal path to fail validation")
	}
	if !strings.Contains(err.Error(), "must be inside outputRoot") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateCleanPathRejectsFileTarget(t *testing.T) {
	outputRoot := ioCleanupFixturePath(t, "file_target")
	filePath := filepath.Join(outputRoot, "data", "maf", "function", "generated", "single.mcfunction")

	err := validateCleanPath(outputRoot, filePath)
	if err == nil {
		t.Fatal("expected file target to fail validation")
	}
	if !strings.Contains(err.Error(), "target must be directory") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateCleanPathRejectsUnsupportedExtensionRecursively(t *testing.T) {
	outputRoot := ioCleanupFixturePath(t, "unsafe_tree")
	cleanPath := filepath.Join(outputRoot, "data", "maf", "function", "generated")

	err := validateCleanPath(outputRoot, cleanPath)
	if err == nil {
		t.Fatal("expected unsupported extension to fail validation")
	}
	if !strings.Contains(err.Error(), "unsupported extension") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCleanupConfiguredPathsRemovesValidatedPaths(t *testing.T) {
	outputRoot := ioCleanupFixturePath(t, "cleanup_tree")
	first := filepath.Join(outputRoot, "data", "maf", "function", "generated")
	second := filepath.Join(outputRoot, "data", "maf", "loot_table", "generated")

	var removed []string
	mockRemoveAllPath(t, func(path string) error {
		removed = append(removed, path)
		return nil
	})

	if err := cleanupConfiguredPaths(outputRoot, []string{first, second}); err != nil {
		t.Fatalf("cleanup should succeed: %v", err)
	}
	if !reflect.DeepEqual(removed, []string{first, second}) {
		t.Fatalf("unexpected remove paths: got=%v", removed)
	}
}

func TestCleanupConfiguredPathsValidatesAllBeforeRemoving(t *testing.T) {
	outputRoot := ioCleanupFixturePath(t, "cleanup_tree_invalid")
	validPath := filepath.Join(outputRoot, "data", "maf", "function", "generated")
	invalidPath := filepath.Join(outputRoot, "data", "maf", "loot_table", "generated")

	calls := 0
	mockRemoveAllPath(t, func(path string) error {
		calls++
		return nil
	})

	err := cleanupConfiguredPaths(outputRoot, []string{validPath, invalidPath})
	if err == nil {
		t.Fatal("expected cleanup to fail when one path is invalid")
	}
	if calls != 0 {
		t.Fatalf("removeAll must not be called when validation fails, calls=%d", calls)
	}
}

func TestCleanupConfiguredPathsPropagatesRemoveError(t *testing.T) {
	outputRoot := ioCleanupFixturePath(t, "cleanup_tree")
	first := filepath.Join(outputRoot, "data", "maf", "function", "generated")
	second := filepath.Join(outputRoot, "data", "maf", "loot_table", "generated")

	mockRemoveAllPath(t, func(path string) error {
		if path == second {
			return errors.New("boom")
		}
		return nil
	})

	err := cleanupConfiguredPaths(outputRoot, []string{first, second})
	if err == nil {
		t.Fatal("expected cleanup to fail on removeAll error")
	}
	if !strings.Contains(err.Error(), "cleanPaths[1]") {
		t.Fatalf("error should include failing index: %v", err)
	}
}

func TestCleanupConfiguredPathsNoopWhenEmpty(t *testing.T) {
	outputRoot := ioCleanupFixturePath(t, "noop_tree")

	calls := 0
	mockRemoveAllPath(t, func(path string) error {
		calls++
		return nil
	})

	if err := cleanupConfiguredPaths(outputRoot, nil); err != nil {
		t.Fatalf("empty cleanPaths should be noop: %v", err)
	}
	if calls != 0 {
		t.Fatalf("removeAll must not be called for empty cleanPaths, calls=%d", calls)
	}
}

func mockRemoveAllPath(t *testing.T, fn func(path string) error) {
	t.Helper()
	prev := removeAllPath
	removeAllPath = fn
	t.Cleanup(func() {
		removeAllPath = prev
	})
}

func ioCleanupFixturePath(t *testing.T, name string) string {
	t.Helper()
	return filepath.Join(mustGetwd(t), "testdata", "io_cleanup", name)
}

func mustGetwd(t *testing.T) string {
	t.Helper()
	wd, err := filepath.Abs(".")
	if err != nil {
		t.Fatal(err)
	}
	return wd
}
