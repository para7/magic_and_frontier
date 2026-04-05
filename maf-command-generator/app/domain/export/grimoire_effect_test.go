package export

import (
	"testing"

	grimoireModel "maf_command_editor/app/domain/model/grimoire"
)

func TestBuildGrimoireArtifactsBuildsEffectsAndSelectExec(t *testing.T) {
	master := exportMasterStub{
		grimoires: []grimoireModel.Grimoire{
			{ID: "fire", Script: []string{"say fire"}},
			{ID: "ice", Script: []string{"say ice"}},
		},
	}

	effects := BuildGrimoireArtifacts(master)

	if len(effects) != 2 {
		t.Fatalf("effects length = %d, want 2", len(effects))
	}
	if effects[0].ID != "fire" || effects[0].Body != "say fire" {
		t.Fatalf("effects[0] = %#v", effects[0])
	}
	if effects[1].ID != "ice" || effects[1].Body != "say ice" {
		t.Fatalf("effects[1] = %#v", effects[1])
	}

	if effects[0].Book == "" {
		t.Fatalf("effects[0].Book should be populated")
	}
	if effects[1].Book == "" {
		t.Fatalf("effects[1].Book should be populated")
	}
}

func TestBuildGrimoireArtifactsEmpty(t *testing.T) {
	effects := BuildGrimoireArtifacts(exportMasterStub{})

	if len(effects) != 0 {
		t.Fatalf("effects length = %d, want 0", len(effects))
	}
}
