package web

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"tools2/app/internal/webui"
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

	unbreaking := findEnchantmentOption(t, options, "minecraft:unbreaking")
	if unbreaking.Checked || unbreaking.Level != "3" {
		t.Fatalf("unbreaking = %#v", unbreaking)
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
	if enchantments != "minecraft:sharpness 5\nminecraft:unbreaking 3" {
		t.Fatalf("enchantments = %q", enchantments)
	}

	sharpness := findEnchantmentOption(t, options, "minecraft:sharpness")
	if !sharpness.Checked || sharpness.Level != "5" {
		t.Fatalf("sharpness = %#v", sharpness)
	}

	mending := findEnchantmentOption(t, options, "minecraft:mending")
	if mending.Checked || mending.Level != "7" {
		t.Fatalf("mending = %#v", mending)
	}
	if mending.LevelFieldName != "enchantmentLevel.mending" {
		t.Fatalf("mending level field name = %q", mending.LevelFieldName)
	}
}

func findEnchantmentOption(t *testing.T, options []webui.ItemEnchantmentOption, id string) webui.ItemEnchantmentOption {
	t.Helper()
	for _, option := range options {
		if option.ID == id {
			return option
		}
	}
	t.Fatalf("enchantment option %q not found", id)
	return webui.ItemEnchantmentOption{}
}
