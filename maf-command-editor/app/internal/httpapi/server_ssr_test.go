package httpapi

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"tools2/app/internal/domain/grimoire"
	"tools2/app/internal/domain/items"
	"tools2/app/internal/domain/skills"
)

func TestHandlerHealth(t *testing.T) {
	handler, _ := newTestHandler(t)
	rec := request(t, handler, http.MethodGet, "/health", nil, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	var body map[string]bool
	decodeResponse(t, rec, &body)
	if !body["ok"] {
		t.Fatalf("body = %#v", body)
	}
}

func TestHandlerSSRSkillCreate(t *testing.T) {
	handler, _ := newTestHandler(t)

	rec := request(t, handler, http.MethodGet, "/skills/new", nil, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `action="/skills/new"`) {
		t.Fatalf("body = %s", rec.Body.String())
	}

	rec = postForm(t, handler, "/skills/new", url.Values{
		"id":          {"skill_1"},
		"name":        {"Slash"},
		"description": {"Basic skill"},
		"script":      {"say slash"},
	}, http.StatusSeeOther)
	assertRedirect(t, rec, "/skills")

	rec = request(t, handler, http.MethodGet, "/skills", nil, "")
	if !strings.Contains(rec.Body.String(), "Slash") {
		t.Fatalf("body = %s", rec.Body.String())
	}
}

func TestHandlerSSRGrimoireEditShowsReadonlyCastID(t *testing.T) {
	handler, _ := newTestHandler(t)
	grimoireID := createJSONEntry(t, handler, "/api/grimoire", grimoire.SaveInput{
		ID:       "grimoire_1",
		Title:    "Firebolt",
		Script:   "say fire",
		CastTime: 20,
		MPCost:   5,
	})

	rec := request(t, handler, http.MethodGet, "/grimoire/edit?id="+grimoireID, nil, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `name="castid" value="1" readonly`) {
		t.Fatalf("body = %s", rec.Body.String())
	}
}

func TestHandlerSSRTreasureCreateValidationKeepsIDEditable(t *testing.T) {
	handler, _ := newTestHandler(t)

	rec := postForm(t, handler, "/treasures/edit", url.Values{
		"id":            {"bad"},
		"tablePath":     {"minecraft:chests/missing_entry"},
		"lootPoolsText": {"minecraft_item,minecraft:apple,1,,"},
	}, http.StatusOK)
	body := rec.Body.String()
	if !strings.Contains(body, `name="id" value="bad"`) {
		t.Fatalf("body = %s", body)
	}
	if strings.Contains(body, `name="id" value="bad" readonly`) {
		t.Fatalf("body = %s", body)
	}
}

func TestHandlerSSRItemsListIncludesClientControlsAndReturnTo(t *testing.T) {
	handler, _ := newTestHandler(t)
	skillID := createJSONEntry(t, handler, "/api/skills", skills.SaveInput{
		ID:     "skill_1",
		Name:   "Slash",
		Script: "say slash",
	})
	itemID := createJSONEntry(t, handler, "/api/items", items.SaveInput{
		ID:      "items_1",
		ItemID:  "minecraft:apple",
		SkillID: skillID,
	})

	rec := request(t, handler, http.MethodGet, "/items?q=apple&page=2", nil, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, `data-list-root`) || !strings.Contains(body, `data-list-search`) {
		t.Fatalf("body = %s", body)
	}
	if !strings.Contains(body, `href="/items/new?returnTo=%2Fitems%3Fq%3Dapple%26page%3D2"`) {
		t.Fatalf("body = %s", body)
	}
	wantEdit := `href="/items/edit?id=` + itemID + `&amp;returnTo=%2Fitems%3Fq%3Dapple%26page%3D2"`
	if !strings.Contains(body, wantEdit) {
		t.Fatalf("body = %s", body)
	}
	if !strings.Contains(body, `data-sort-item_id="minecraft:apple"`) {
		t.Fatalf("body = %s", body)
	}
}

func TestHandlerSSRIncludesDeleteConfirmScript(t *testing.T) {
	handler, _ := newTestHandler(t)

	rec := request(t, handler, http.MethodGet, "/items", nil, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, `document.addEventListener("submit", (event) => {`) {
		t.Fatalf("body = %s", body)
	}
	if !strings.Contains(body, `actionURL.pathname.endsWith("/delete")`) {
		t.Fatalf("body = %s", body)
	}
	if !strings.Contains(body, `window.confirm(deleteConfirmMessage)`) {
		t.Fatalf("body = %s", body)
	}
	if !strings.Contains(body, `削除してもよろしいですか？`) {
		t.Fatalf("body = %s", body)
	}
}

func TestHandlerSSRSkillEditRespectsReturnToOnSaveAndFallback(t *testing.T) {
	handler, _ := newTestHandler(t)
	skillID := createJSONEntry(t, handler, "/api/skills", skills.SaveInput{
		ID:          "skill_1",
		Name:        "Slash",
		Description: "Basic slash",
		Script:      "say slash",
	})

	rec := request(t, handler, http.MethodGet, "/skills/edit?id="+skillID+"&returnTo=%2Fskills%3Fq%3Dslash%26page%3D2", nil, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `href="/skills?q=slash&amp;page=2">キャンセル</a>`) {
		t.Fatalf("body = %s", rec.Body.String())
	}

	rec = postForm(t, handler, "/skills/edit", url.Values{
		"id":          {skillID},
		"name":        {"Slash v2"},
		"description": {"Updated"},
		"script":      {"say slash2"},
		"returnTo":    {"/skills?q=slash&page=2"},
	}, http.StatusSeeOther)
	assertRedirect(t, rec, "/skills?q=slash&page=2")

	rec = postForm(t, handler, "/skills/edit", url.Values{
		"id":          {skillID},
		"name":        {"Slash v3"},
		"description": {"Updated again"},
		"script":      {"say slash3"},
		"returnTo":    {"https://evil.example"},
	}, http.StatusSeeOther)
	assertRedirect(t, rec, "/skills")
}

func TestHandlerSSRItemReturnToRejectsNonListPaths(t *testing.T) {
	handler, _ := newTestHandler(t)
	itemID := createJSONEntry(t, handler, "/api/items", items.SaveInput{
		ID:     "items_1",
		ItemID: "minecraft:apple",
	})

	rec := request(t, handler, http.MethodGet, "/items/edit?id="+itemID+"&returnTo=%2Fsave", nil, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `href="/items">キャンセル</a>`) {
		t.Fatalf("body = %s", rec.Body.String())
	}

	rec = postForm(t, handler, "/items/edit", url.Values{
		"id":       {itemID},
		"itemId":   {"minecraft:apple"},
		"returnTo": {"/items/edit?id=" + itemID},
	}, http.StatusSeeOther)
	assertRedirect(t, rec, "/items")
}

func TestHandlerSSRItemFormShowsEnchantmentControls(t *testing.T) {
	handler, _ := newTestHandler(t)

	rec := request(t, handler, http.MethodGet, "/items/new", nil, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, `class="enchantments-details"`) {
		t.Fatalf("body = %s", body)
	}
	if !strings.Contains(body, `name="enchantmentIds" value="minecraft:sharpness"`) {
		t.Fatalf("body = %s", body)
	}
	if strings.Contains(body, `class="code enchantment-id"`) {
		t.Fatalf("body = %s", body)
	}
	if !strings.Contains(body, `name="enchantmentLevel.sharpness"`) {
		t.Fatalf("body = %s", body)
	}
	if !strings.Contains(body, `class="enchantments-category-title">Armor</h4>`) {
		t.Fatalf("body = %s", body)
	}
	if !strings.Contains(body, `class="enchantments-category-title">Weapon</h4>`) {
		t.Fatalf("body = %s", body)
	}
	if !strings.Contains(body, `class="enchantments-category-title">Durability</h4>`) {
		t.Fatalf("body = %s", body)
	}

	armorAt := strings.Index(body, `class="enchantments-category-title">Armor</h4>`)
	weaponAt := strings.Index(body, `class="enchantments-category-title">Weapon</h4>`)
	durabilityAt := strings.Index(body, `class="enchantments-category-title">Durability</h4>`)
	if !(armorAt < weaponAt && weaponAt < durabilityAt) {
		t.Fatalf("unexpected enchantment category order: armor=%d weapon=%d durability=%d", armorAt, weaponAt, durabilityAt)
	}
}

func TestHandlerSSRItemFormOpensEnchantmentsOnValidationError(t *testing.T) {
	handler, _ := newTestHandler(t)

	rec := postForm(t, handler, "/items/new", url.Values{
		"id":                         {"items_1"},
		"itemId":                     {"minecraft:stone"},
		"enchantmentIds":             {"minecraft:sharpness"},
		"enchantmentLevel.sharpness": {"abc"},
	}, http.StatusOK)

	body := rec.Body.String()
	if !strings.Contains(body, `class="enchantments-details" open`) {
		t.Fatalf("body = %s", body)
	}
	if !strings.Contains(body, "Invalid enchantment line:") {
		t.Fatalf("body = %s", body)
	}
}
