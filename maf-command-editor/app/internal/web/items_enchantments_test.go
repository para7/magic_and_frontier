package web

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"tools2/app/internal/web/ui"
)

func TestItemFormEnchantmentsFromText(t *testing.T) {
	options, selected := itemFormEnchantmentsFromText("minecraft:sharpness 3\nminecraft:mending 1")
	if selected != 2 {
		t.Fatalf("selected = %d", selected)
	}

	sharpness := findEnchantmentOption(t, options, "minecraft:sharpness")
	if !sharpness.Checked || sharpness.Level != "3" {
		t.Fatalf("sharpness = %#v", sharpness)
	}

	mending := findEnchantmentOption(t, options, "minecraft:mending")
	if !mending.Checked || mending.Level != "1" {
		t.Fatalf("mending = %#v", mending)
	}
	if mending.Category != "Durability" {
		t.Fatalf("mending category = %q", mending.Category)
	}

	unbreaking := findEnchantmentOption(t, options, "minecraft:unbreaking")
	if unbreaking.Checked || unbreaking.Level != "0" {
		t.Fatalf("unbreaking = %#v", unbreaking)
	}
}

func TestItemFormEnchantmentsCategoryOrder(t *testing.T) {
	options, _ := itemFormEnchantmentsFromText("")

	want := []string{
		"Armor",
		"Movement/Environment",
		"Weapon",
		"Bow",
		"Crossbow",
		"Bow/Crossbow",
		"Trident",
		"Tools",
		"Durability",
		"Fishing",
		"Cursed",
	}

	got := make([]string, 0, len(want))
	seen := map[string]bool{}
	for _, option := range options {
		if option.Category == "" || seen[option.Category] {
			continue
		}
		seen[option.Category] = true
		got = append(got, option.Category)
	}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("categories = %#v", got)
	}
}

func TestItemFormEnchantmentsFromRequest(t *testing.T) {
	values := url.Values{
		"enchantmentIds":                 {"minecraft:sharpness", "minecraft:unbreaking"},
		"enchantmentLevel.sharpness":     {"5"},
		"enchantmentLevel.mending":       {"7"},
		"enchantmentLevel.unbreaking":    {""},
		"enchantmentLevel.aqua_affinity": {"1"},
	}
	req := httptest.NewRequest(http.MethodPost, "/items/new", strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err := req.ParseForm(); err != nil {
		t.Fatal(err)
	}

	enchantments, options, selected := itemFormEnchantmentsFromRequest(req)
	if selected != 2 {
		t.Fatalf("selected = %d", selected)
	}
	if enchantments != "minecraft:sharpness 5\nminecraft:unbreaking 0" {
		t.Fatalf("enchantments = %q", enchantments)
	}

	sharpness := findEnchantmentOption(t, options, "minecraft:sharpness")
	if !sharpness.Checked || sharpness.Level != "5" {
		t.Fatalf("sharpness = %#v", sharpness)
	}
	if sharpness.Category != "Weapon" {
		t.Fatalf("sharpness category = %q", sharpness.Category)
	}

	mending := findEnchantmentOption(t, options, "minecraft:mending")
	if mending.Checked || mending.Level != "7" {
		t.Fatalf("mending = %#v", mending)
	}
	if mending.LevelFieldName != "enchantmentLevel.mending" {
		t.Fatalf("mending level field name = %q", mending.LevelFieldName)
	}
}

func findEnchantmentOption(t *testing.T, options []ui.ItemEnchantmentOption, id string) ui.ItemEnchantmentOption {
	t.Helper()
	for _, option := range options {
		if option.ID == id {
			return option
		}
	}
	t.Fatalf("enchantment option %q not found", id)
	return ui.ItemEnchantmentOption{}
}
