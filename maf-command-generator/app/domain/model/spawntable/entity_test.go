package spawntable

import (
	"os"
	"path/filepath"
	"testing"

	model "maf_command_editor/app/domain/model"
)

func validSpawnTable() SpawnTable {
	return SpawnTable{
		ID:            "spawn_1",
		SourceMobType: "minecraft:zombie",
		Dimension:     "minecraft:overworld",
		MinDistance:   0,
		MaxDistance:   1000,
		MinX:          -100, MaxX: 100,
		MinY: -64, MaxY: 320,
		MinZ: -100, MaxZ: 100,
		BaseMob: &BaseMob{Weight: 10},
		Replacements: []model.ReplacementEntry{
			{EnemyID: "enemy_1", Weight: 5},
		},
	}
}

func floatPtr(v float64) *float64 {
	return &v
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
		{name: "minDistance ng out of range", patch: func(st *SpawnTable) { st.MinDistance = -1 }, wantErrField: "minDistance"},
		{name: "maxDistance ng out of range", patch: func(st *SpawnTable) { st.MaxDistance = 100000000 }, wantErrField: "maxDistance"},
		{name: "minX ok boundary", patch: func(st *SpawnTable) { st.MinX = -99999999 }},
		{name: "minX ng out of range", patch: func(st *SpawnTable) { st.MinX = -100000000 }, wantErrField: "minX"},
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

func TestSpawnTableValidateRelationBaseMobAttributesRange(t *testing.T) {
	entity := &SpawnTableEntity{}
	tests := []struct {
		name         string
		patch        func(*SpawnTable)
		wantErrField string
	}{
		{
			name: "hp too low",
			patch: func(st *SpawnTable) {
				st.BaseMob.Attributes = &BaseMobAttributes{HP: floatPtr(0)}
			},
			wantErrField: "baseMob.attributes.hp",
		},
		{
			name: "attack too low",
			patch: func(st *SpawnTable) {
				st.BaseMob.Attributes = &BaseMobAttributes{Attack: floatPtr(-1)}
			},
			wantErrField: "baseMob.attributes.attack",
		},
		{
			name: "defense too high",
			patch: func(st *SpawnTable) {
				st.BaseMob.Attributes = &BaseMobAttributes{Defense: floatPtr(100001)}
			},
			wantErrField: "baseMob.attributes.defense",
		},
		{
			name: "movespeed valid",
			patch: func(st *SpawnTable) {
				st.BaseMob.Attributes = &BaseMobAttributes{MoveSpeed: floatPtr(0.22)}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := validSpawnTable()
			tt.patch(&st)
			errs := entity.ValidateRelation(st, testDBMaster{})
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
		MinDistance:   0,
		MaxDistance:   100,
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

	nonDistance := base
	nonDistance.MinDistance = 200
	nonDistance.MaxDistance = 300
	if len(AllOverlaps([]SpawnTable{base, nonDistance})) != 0 {
		t.Fatal("expected non-overlapping distance ranges to not be detected")
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

func TestSpawnTableLoadCoordinatesFormat(t *testing.T) {
	dir := t.TempDir()
	jsonPath := filepath.Join(dir, "over1.json")
	raw := `{
  "coordinates": {
    "dimension": "minecraft:overworld",
    "minDistance": 0,
    "maxDistance": 1000,
    "minX": -99999999,
    "maxX": 99999999,
    "minZ": -99999999,
    "maxZ": 99999999,
    "minY": -64,
    "maxY": 320
  },
  "entries": [
    {
      "sourceMobType": "minecraft:zombie",
      "baseMob": {
        "weight": 30,
        "attributes": {
          "hp": 100
        }
      },
      "replacements": [
        {"enemyId": "enemy_1", "weight": 70}
      ]
    }
  ]
}`
	if err := os.WriteFile(jsonPath, []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}

	entity := NewSpawnTableEntity(dir)
	if err := entity.Load(); err != nil {
		t.Fatal(err)
	}
	got := entity.GetAll()
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
	if got[0].ID != "over1_zombie_1" {
		t.Fatalf("unexpected generated id: %q", got[0].ID)
	}
	if got[0].MinDistance != 0 || got[0].MaxDistance != 1000 {
		t.Fatalf("unexpected distance range: %+v", got[0])
	}
	if got[0].Dimension != "minecraft:overworld" {
		t.Fatalf("unexpected dimension: %q", got[0].Dimension)
	}
	if got[0].GetBaseMobWeight() != 30 {
		t.Fatalf("unexpected base mob weight: %+v", got[0].BaseMob)
	}
	if got[0].BaseMob == nil || got[0].BaseMob.Attributes == nil || got[0].BaseMob.Attributes.HP == nil || *got[0].BaseMob.Attributes.HP != 100 {
		t.Fatalf("unexpected base mob attributes: %+v", got[0].BaseMob)
	}
}

func TestSpawnTableLoadLegacyBaseMobWeight(t *testing.T) {
	dir := t.TempDir()
	jsonPath := filepath.Join(dir, "legacy.json")
	raw := `{
  "entries": [
    {
      "id": "legacy_spawn",
      "sourceMobType": "minecraft:zombie",
      "dimension": "minecraft:overworld",
      "minDistance": 0,
      "maxDistance": 1000,
      "minX": -100,
      "maxX": 100,
      "minY": -64,
      "maxY": 320,
      "minZ": -100,
      "maxZ": 100,
      "baseMobWeight": 30,
      "replacements": [
        {"enemyId": "enemy_1", "weight": 70}
      ]
    }
  ]
}`
	if err := os.WriteFile(jsonPath, []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}

	entity := NewSpawnTableEntity(dir)
	if err := entity.Load(); err != nil {
		t.Fatal(err)
	}
	got := entity.GetAll()
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
	if got[0].GetBaseMobWeight() != 30 {
		t.Fatalf("unexpected base mob weight: %+v", got[0].BaseMob)
	}
	if got[0].BaseMob == nil || got[0].BaseMob.Attributes != nil {
		t.Fatalf("unexpected base mob for legacy format: %+v", got[0].BaseMob)
	}
}
