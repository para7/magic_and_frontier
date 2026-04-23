package enemy

import (
	"encoding/json"
	"testing"
)

func TestEnemyUnmarshalJSONPreservesMinecraftNumbers(t *testing.T) {
	raw := []byte(`{
		"id":"enemy_1",
		"mobType":"minecraft:zombie",
		"minecraft":{
			"Health":2,
			"attributes":[{"id":"minecraft:movement_speed","base":0.22}]
		},
		"maf":{"dropMode":"replace","drops":[]}
	}`)

	var got Enemy
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal enemy: %v", err)
	}

	health, ok := got.Minecraft["Health"].(json.Number)
	if !ok || health.String() != "2" {
		t.Fatalf("expected minecraft.Health as json.Number(2), got %#v", got.Minecraft["Health"])
	}

	attrs, ok := got.Minecraft["attributes"].([]any)
	if !ok || len(attrs) != 1 {
		t.Fatalf("unexpected attributes: %#v", got.Minecraft["attributes"])
	}
	attr0, ok := attrs[0].(map[string]any)
	if !ok {
		t.Fatalf("unexpected attribute[0]: %#v", attrs[0])
	}
	base, ok := attr0["base"].(json.Number)
	if !ok || base.String() != "0.22" {
		t.Fatalf("expected attributes[0].base as json.Number(0.22), got %#v", attr0["base"])
	}
}
