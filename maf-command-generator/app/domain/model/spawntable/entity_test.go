package spawntable

import (
	"testing"

	model "maf_command_editor/app/domain/model"
)

func validSpawnTable() SpawnTable {
	return SpawnTable{
		ID:            "spawn_1",
		SourceMobType: "minecraft:zombie",
		Dimension:     "minecraft:overworld",
		MinX:          -100, MaxX: 100,
		MinY: -64, MaxY: 320,
		MinZ: -100, MaxZ: 100,
		BaseMobWeight: 10,
		Replacements: []model.ReplacementEntry{
			{EnemyID: "enemy_1", Weight: 5},
		},
	}
}

func hasFieldError(errs []model.ValidationError, field string) bool {
	for _, err := range errs {
		if err.Field == field {
			return true
		}
	}
	return false
}

type testDBMaster struct{}

func (testDBMaster) HasItem(string) bool     { return true }
func (testDBMaster) HasGrimoire(string) bool { return true }
func (testDBMaster) HasPassive(string) bool  { return true }
func (testDBMaster) GetPassive(string) (model.PassiveSnapshot, bool) {
	v := true
	return model.PassiveSnapshot{ID: "passive_1", GenerateGrimoire: &v}, true
}
func (testDBMaster) HasBow(string) bool                { return true }
func (testDBMaster) HasEnemySkill(string) bool         { return true }
func (testDBMaster) HasEnemy(string) bool              { return true }
func (testDBMaster) HasSpawnTable(string) bool         { return true }
func (testDBMaster) HasTreasure(string) bool           { return true }
func (testDBMaster) HasMinecraftLootTable(string) bool { return true }

func TestSpawnTableValidateStructAllValid(t *testing.T) {
	entity := &SpawnTableEntity{}
	errs := entity.ValidateStruct(validSpawnTable())
	if len(errs) != 0 {
		t.Fatalf("expected no validation errors, got %#v", errs)
	}
}

func TestSpawnTableValidateStructPerField(t *testing.T) {
	entity := &SpawnTableEntity{}

	tests := []struct {
		name         string
		patch        func(*SpawnTable)
		wantErrField string
	}{
		{name: "id ok", patch: func(st *SpawnTable) { st.ID = "ok" }},
		{name: "id ng empty", patch: func(st *SpawnTable) { st.ID = "" }, wantErrField: "id"},
		{name: "sourceMobType ok", patch: func(st *SpawnTable) { st.SourceMobType = "minecraft:skeleton" }},
		{name: "sourceMobType ng empty", patch: func(st *SpawnTable) { st.SourceMobType = " " }, wantErrField: "sourceMobType"},
		{name: "dimension ok overworld", patch: func(st *SpawnTable) { st.Dimension = "minecraft:overworld" }},
		{name: "dimension ok nether", patch: func(st *SpawnTable) { st.Dimension = "minecraft:the_nether" }},
		{name: "dimension ok end", patch: func(st *SpawnTable) { st.Dimension = "minecraft:the_end" }},
		{name: "dimension ng invalid", patch: func(st *SpawnTable) { st.Dimension = "minecraft:invalid" }, wantErrField: "dimension"},
		{name: "minX ok boundary", patch: func(st *SpawnTable) { st.MinX = -30000000 }},
		{name: "minX ng out of range", patch: func(st *SpawnTable) { st.MinX = -30000001 }, wantErrField: "minX"},
		{name: "replacements ng empty", patch: func(st *SpawnTable) { st.Replacements = nil }, wantErrField: "replacements"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := validSpawnTable()
			tt.patch(&st)
			errs := entity.ValidateStruct(st)
			if tt.wantErrField == "" {
				if len(errs) != 0 {
					t.Fatalf("expected no errors, got %#v", errs)
				}
				return
			}
			if !hasFieldError(errs, tt.wantErrField) {
				t.Fatalf("expected error for field %q, got %#v", tt.wantErrField, errs)
			}
		})
	}
}

func TestAllOverlaps(t *testing.T) {
	base := SpawnTable{
		SourceMobType: "minecraft:zombie",
		Dimension:     "minecraft:overworld",
		MinX:          0, MaxX: 100, MinY: 0, MaxY: 100, MinZ: 0, MaxZ: 100,
	}
	other := base
	other.MinX = 50
	other.MaxX = 150
	if len(AllOverlaps([]SpawnTable{base, other})) == 0 {
		t.Fatal("expected overlapping spawn tables to be detected")
	}

	nonOverlap := base
	nonOverlap.MinX = 200
	nonOverlap.MaxX = 300
	if len(AllOverlaps([]SpawnTable{base, nonOverlap})) != 0 {
		t.Fatal("expected non-overlapping spawn tables to not be detected")
	}
}

func TestSpawnTableValidateAllDetectsDuplicateID(t *testing.T) {
	entity := &SpawnTableEntity{
		data: []SpawnTable{
			validSpawnTable(),
			validSpawnTable(),
		},
	}

	allErrs := entity.ValidateAll(testDBMaster{})
	for _, recordErrs := range allErrs {
		for _, err := range recordErrs {
			if err.ID == "spawn_1" && err.Field == "id" && err.Tag == "unique" {
				return
			}
		}
	}
	t.Fatalf("expected duplicate id error, got %#v", allErrs)
}
