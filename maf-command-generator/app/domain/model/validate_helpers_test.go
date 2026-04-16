package model

import (
	"strings"
	"testing"
)

type validateHelpersMasterStub struct {
	passive bool
}

func (s validateHelpersMasterStub) HasItem(string) bool     { return true }
func (s validateHelpersMasterStub) HasGrimoire(string) bool { return true }
func (s validateHelpersMasterStub) HasPassive(string) bool  { return s.passive }
func (s validateHelpersMasterStub) GetPassive(string) (PassiveSnapshot, bool) {
	v := true
	return PassiveSnapshot{ID: "passive_1", GenerateGrimoire: &v}, true
}
func (s validateHelpersMasterStub) HasBow(string) bool                { return true }
func (s validateHelpersMasterStub) HasEnemySkill(string) bool         { return true }
func (s validateHelpersMasterStub) HasEnemy(string) bool              { return true }
func (s validateHelpersMasterStub) HasSpawnTable(string) bool         { return true }
func (s validateHelpersMasterStub) HasTreasure(string) bool           { return true }
func (s validateHelpersMasterStub) HasMinecraftLootTable(string) bool { return true }

func TestValidateDropRefsPassiveRequiresSlot(t *testing.T) {
	errs := ValidateDropRefs("enemy", "enemy_1", "drops", []DropRef{
		{Kind: "passive", RefID: "passive_1", Weight: 1},
	}, validateHelpersMasterStub{passive: true})
	if len(errs) != 1 || errs[0].Field != "drops[0].slot" {
		t.Fatalf("expected slot relation error, got %#v", errs)
	}
}

func TestValidateDropRefsPassiveChecksRelation(t *testing.T) {
	slot := 1
	errs := ValidateDropRefs("enemy", "enemy_1", "drops", []DropRef{
		{Kind: "passive", RefID: "passive_missing", Slot: &slot, Weight: 1},
	}, validateHelpersMasterStub{passive: false})
	if len(errs) != 1 || errs[0].Field != "drops[0].refId" {
		t.Fatalf("expected refId relation error, got %#v", errs)
	}
}

func TestValidateDropRefsNonPassiveRejectsSlot(t *testing.T) {
	slot := 2
	errs := ValidateDropRefs("enemy", "enemy_1", "drops", []DropRef{
		{Kind: "item", RefID: "item_1", Slot: &slot, Weight: 1},
	}, validateHelpersMasterStub{passive: true})
	if len(errs) != 1 || errs[0].Field != "drops[0].slot" {
		t.Fatalf("expected slot unsupported error, got %#v", errs)
	}
}

func TestValidateMafLootPoolsPassiveRequiresSlot(t *testing.T) {
	errs := ValidateMafLootPools("enemy", "enemy_1", "drops", []any{
		map[string]any{
			"entries": []any{
				map[string]any{
					"type": "maf:passive",
					"name": "passive_1",
				},
			},
		},
	}, validateHelpersMasterStub{passive: true})
	if len(errs) != 1 || errs[0].Field != "drops[0].entries[0].slot" {
		t.Fatalf("expected slot required error, got %#v", errs)
	}
}

func TestValidateMafLootPoolsRejectsInvertedCountRange(t *testing.T) {
	errs := ValidateMafLootPools("enemy", "enemy_1", "drops", []any{
		map[string]any{
			"entries": []any{
				map[string]any{
					"type": "maf:item",
					"name": "item_1",
					"count": map[string]any{
						"min": 3.0,
						"max": 1.0,
					},
				},
			},
		},
	}, validateHelpersMasterStub{passive: true})
	if len(errs) != 1 {
		t.Fatalf("expected exactly one error, got %#v", errs)
	}
	if errs[0].Field != "drops[0].entries[0].count" {
		t.Fatalf("unexpected field: %#v", errs[0])
	}
	if !strings.Contains(errs[0].Param, "less than or equal") {
		t.Fatalf("unexpected error param: %#v", errs[0])
	}
}

func TestValidateMafLootPoolsRejectsUnsupportedMafType(t *testing.T) {
	errs := ValidateMafLootPools("enemy", "enemy_1", "drops", []any{
		map[string]any{
			"entries": []any{
				map[string]any{
					"type": "maf:unknown",
					"name": "x",
				},
			},
		},
	}, validateHelpersMasterStub{passive: true})
	if len(errs) != 1 || errs[0].Field != "drops[0].entries[0].type" {
		t.Fatalf("expected unsupported type error, got %#v", errs)
	}
}
