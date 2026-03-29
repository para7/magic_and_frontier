package export

import "testing"

func TestLootTableOutputPathRejectsTraversal(t *testing.T) {
	settings := ExportSettings{OutputRoot: "/tmp/out"}
	if _, err := lootTableOutputPath(settings, "maf:loot/../escape"); err == nil {
		t.Fatalf("expected traversal table path to be rejected")
	}
}
