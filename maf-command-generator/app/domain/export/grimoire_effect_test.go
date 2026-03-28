package export

// import (
// 	"testing"

// 	grimoireModel "maf_command_editor/app/domain/model/grimoire"
// )

// type stubDBMaster struct {
// 	entries []grimoireModel.Grimoire
// }

// func (s stubDBMaster) GetGrimoireByID(id string) (grimoireModel.Grimoire, bool) {
// 	for _, entry := range s.entries {
// 		if entry.ID == id {
// 			return entry, true
// 		}
// 	}
// 	return grimoireModel.Grimoire{}, false
// }

// func (s stubDBMaster) ListGrimoires() []grimoireModel.Grimoire {
// 	out := make([]grimoireModel.Grimoire, len(s.entries))
// 	copy(out, s.entries)
// 	return out
// }

// func TestBuildGrimoireEffectArtifactsBuildsEffectsAndSelectExec(t *testing.T) {
// 	master := stubDBMaster{
// 		entries: []grimoireModel.Grimoire{
// 			{ID: "fire", CastID: 2, Script: "say fire"},
// 			{ID: "ice", CastID: 9, Script: "say ice\n"},
// 		},
// 	}

// 	artifacts := BuildGrimoireEffectArtifacts(master, "maf", "data/maf/function/generated/grimoire")

// 	if len(artifacts.Effects) != 2 {
// 		t.Fatalf("Effects length = %d, want 2", len(artifacts.Effects))
// 	}
// 	if artifacts.Effects[0].ID != "fire" || artifacts.Effects[0].Body != "say fire\n" {
// 		t.Fatalf("Effects[0] = %#v", artifacts.Effects[0])
// 	}
// 	if artifacts.Effects[1].ID != "ice" || artifacts.Effects[1].Body != "say ice\n" {
// 		t.Fatalf("Effects[1] = %#v", artifacts.Effects[1])
// 	}

// 	wantSelectExec := "" +
// 		"execute if entity @s[scores={mafEffectID=2}] run function maf:generated/grimoire/fire\n" +
// 		"execute if entity @s[scores={mafEffectID=9}] run function maf:generated/grimoire/ice\n"
// 	if artifacts.SelectExec != wantSelectExec {
// 		t.Fatalf("SelectExec = %q, want %q", artifacts.SelectExec, wantSelectExec)
// 	}
// }

// func TestBuildGrimoireEffectArtifactsEmpty(t *testing.T) {
// 	artifacts := BuildGrimoireEffectArtifacts(stubDBMaster{}, "maf", "data/maf/function/generated/grimoire")

// 	if len(artifacts.Effects) != 0 {
// 		t.Fatalf("Effects length = %d, want 0", len(artifacts.Effects))
// 	}
// 	if artifacts.SelectExec != "\n" {
// 		t.Fatalf("SelectExec = %q, want \\n", artifacts.SelectExec)
// 	}
// }
