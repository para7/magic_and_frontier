package export

import (
	"testing"

	grimoireModel "maf_command_editor/app/domain/model/grimoire"
)

func TestBuildGrimoireArtifactsBuildsEffectsAndSelectExec(t *testing.T) {
	master := exportMasterStub{
		grimoires: []grimoireModel.Grimoire{
			{ID: "fire", CastID: 2, Script: []string{"say fire"}},
			{ID: "ice", CastID: 9, Script: []string{"say ice"}},
		},
	}

	effects := BuildGrimoireArtifacts(master, "generated/grimoire/effect")

	if len(effects) != 2 {
		t.Fatalf("effects length = %d, want 2", len(effects))
	}
	if effects[0].ID != "fire" || effects[0].Body != "say fire" {
		t.Fatalf("effects[0] = %#v", effects[0])
	}
	if effects[1].ID != "ice" || effects[1].Body != "say ice" {
		t.Fatalf("effects[1] = %#v", effects[1])
	}

	wantRef0 := "maf:generated/grimoire/effect/fire"
	wantRef1 := "maf:generated/grimoire/effect/ice"
	if effects[0].FunctionRef != wantRef0 {
		t.Fatalf("effects[0].FunctionRef = %q, want %q", effects[0].FunctionRef, wantRef0)
	}
	if effects[1].FunctionRef != wantRef1 {
		t.Fatalf("effects[1].FunctionRef = %q, want %q", effects[1].FunctionRef, wantRef1)
	}
}

func TestBuildGrimoireArtifactsEmpty(t *testing.T) {
	effects := BuildGrimoireArtifacts(exportMasterStub{}, "generated/grimoire/effect")

	if len(effects) != 0 {
		t.Fatalf("effects length = %d, want 0", len(effects))
	}
}
