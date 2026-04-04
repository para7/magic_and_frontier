package model

import "testing"

type validateHelpersMasterStub struct {
	passive bool
}

func (s validateHelpersMasterStub) HasItem(string) bool               { return true }
func (s validateHelpersMasterStub) HasGrimoire(string) bool           { return true }
func (s validateHelpersMasterStub) HasPassive(string) bool            { return s.passive }
func (s validateHelpersMasterStub) HasEnemySkill(string) bool         { return true }
func (s validateHelpersMasterStub) HasEnemy(string) bool              { return true }
func (s validateHelpersMasterStub) HasSpawnTable(string) bool         { return true }
func (s validateHelpersMasterStub) HasTreasure(string) bool           { return true }
func (s validateHelpersMasterStub) HasLootTable(string) bool          { return true }
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
