package export

import (
	"testing"

	grimoireModel "maf_command_editor/app/domain/model/grimoire"
)

type stubDBMaster struct {
	entries []grimoireModel.Grimoire
}

func (s stubDBMaster) GetGrimoireByID(id string) (grimoireModel.Grimoire, bool) {
	for _, entry := range s.entries {
		if entry.ID == id {
			return entry, true
		}
	}
	return grimoireModel.Grimoire{}, false
}

func (s stubDBMaster) ListGrimoires() []grimoireModel.Grimoire {
	out := make([]grimoireModel.Grimoire, len(s.entries))
	copy(out, s.entries)
	return out
}

func TestBuildGrimoireArtifactsBuildsEffectsAndSelectExec(t *testing.T) {
	master := stubDBMaster{
		entries: []grimoireModel.Grimoire{
			{ID: "fire", CastID: 2, Script: "say fire"},
			{ID: "ice", CastID: 9, Script: "say ice\n"},
		},
	}

	effects := BuildGrimoireArtifacts(master, "generated/grimoire/effect")

	if len(effects) != 2 {
		t.Fatalf("effects length = %d, want 2", len(effects))
	}
	if effects[0].ID != "fire" || effects[0].Body != "say fire" {
		t.Fatalf("effects[0] = %#v", effects[0])
	}
	if effects[1].ID != "ice" || effects[1].Body != "say ice\n" {
		t.Fatalf("effects[1] = %#v", effects[1])
	}

	wantSelect0 := "execute if entity @s[scores={mafEffectID=2}] run function maf:generated/grimoire/effect/fire"
	wantSelect1 := "execute if entity @s[scores={mafEffectID=9}] run function maf:generated/grimoire/effect/ice"
	if effects[0].SelectScript != wantSelect0 {
		t.Fatalf("effects[0].SelectScript = %q, want %q", effects[0].SelectScript, wantSelect0)
	}
	if effects[1].SelectScript != wantSelect1 {
		t.Fatalf("effects[1].SelectScript = %q, want %q", effects[1].SelectScript, wantSelect1)
	}
}

func TestBuildGrimoireArtifactsEmpty(t *testing.T) {
	effects := BuildGrimoireArtifacts(stubDBMaster{}, "generated/grimoire/effect")

	if len(effects) != 0 {
		t.Fatalf("effects length = %d, want 0", len(effects))
	}
}
