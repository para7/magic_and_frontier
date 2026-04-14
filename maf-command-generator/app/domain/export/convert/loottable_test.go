package export_convert

import (
	"testing"

	model "maf_command_editor/app/domain/model"
	passiveModel "maf_command_editor/app/domain/model/passive"
)

func TestBuildDropLootPoolAllowsPassiveWithGenerateGrimoireFalse(t *testing.T) {
	falseValue := false
	slot := 1
	_, err := BuildDropLootPool(
		[]model.DropRef{{Kind: "passive", RefID: "passive_1", Slot: &slot, Weight: 1}},
		nil,
		nil,
		map[string]passiveModel.Passive{
			"passive_1": {
				ID:               "passive_1",
				Name:             "Passive",
				Condition:        "always",
				Slots:            []int{1},
				Script:           []string{"say test"},
				GenerateGrimoire: &falseValue,
			},
		},
		nil,
		"enemy(enemy_1)",
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestBuildDropLootPoolPassiveStillRequiresSupportedSlot(t *testing.T) {
	trueValue := true
	slot := 2
	_, err := BuildDropLootPool(
		[]model.DropRef{{Kind: "passive", RefID: "passive_1", Slot: &slot, Weight: 1}},
		nil,
		nil,
		map[string]passiveModel.Passive{
			"passive_1": {
				ID:               "passive_1",
				Name:             "Passive",
				Condition:        "always",
				Slots:            []int{1},
				Script:           []string{"say test"},
				GenerateGrimoire: &trueValue,
			},
		},
		nil,
		"enemy(enemy_1)",
	)
	if err == nil {
		t.Fatal("expected slot validation error")
	}
}
