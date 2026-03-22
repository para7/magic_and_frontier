package treasures

import "testing"

func TestIsSupportedTablePath(t *testing.T) {
	if !IsSupportedTablePath("minecraft:chests/simple_dungeon") {
		t.Fatalf("expected supported path")
	}
	if IsSupportedTablePath("minecraft:blocks/stone") {
		t.Fatalf("expected unsupported path")
	}
}
