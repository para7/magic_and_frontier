package export

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"testing"

	bowModel "maf_command_editor/app/domain/model/bow"
	enemyModel "maf_command_editor/app/domain/model/enemy"
	enemyskillModel "maf_command_editor/app/domain/model/enemyskill"
	grimoireModel "maf_command_editor/app/domain/model/grimoire"
	itemModel "maf_command_editor/app/domain/model/item"
	passiveModel "maf_command_editor/app/domain/model/passive"
	config "maf_command_editor/app/files"
)

var updateGolden = flag.Bool("update", false, "update golden files")

type fixtureEntries[T any] struct {
	Entries []T `json:"entries"`
}

type fixtureCase struct {
	name string
	dir  string
}

func discoverCases(t *testing.T, root string) []fixtureCase {
	t.Helper()

	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatalf("read fixture cases %s: %v", root, err)
	}

	cases := make([]fixtureCase, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		cases = append(cases, fixtureCase{
			name: entry.Name(),
			dir:  filepath.Join(root, entry.Name()),
		})
	}

	slices.SortFunc(cases, func(a, b fixtureCase) int {
		return compareStrings(a.name, b.name)
	})
	return cases
}

func loadFixtureMaster(t *testing.T, caseDir string) exportMasterStub {
	t.Helper()

	inputDir := filepath.Join(caseDir, "input")
	return exportMasterStub{
		grimoires:   loadEntries[grimoireModel.Grimoire](t, inputDir, "grimoires.json"),
		passives:    loadEntries[passiveModel.Passive](t, inputDir, "passives.json"),
		bows:        loadEntries[bowModel.BowPassive](t, inputDir, "bows.json"),
		items:       loadEntries[itemModel.Item](t, inputDir, "items.json"),
		enemySkills: loadEntries[enemyskillModel.EnemySkill](t, inputDir, "enemy_skills.json"),
		enemies:     loadEntries[enemyModel.Enemy](t, inputDir, "enemies.json"),
	}
}

func loadEntries[T any](t *testing.T, inputDir, filename string) []T {
	t.Helper()

	path := filepath.Join(inputDir, filename)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []T{}
	} else if err != nil {
		t.Fatalf("stat fixture input %s: %v", path, err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture input %s: %v", path, err)
	}

	var entries fixtureEntries[T]
	if err := json.Unmarshal(data, &entries); err != nil {
		t.Fatalf("unmarshal fixture input %s: %v", path, err)
	}
	if entries.Entries == nil {
		return []T{}
	}
	return entries.Entries
}

func assertGoldenDir(t *testing.T, goldenDir, actualDir string) {
	t.Helper()

	if *updateGolden {
		if err := syncGoldenDir(goldenDir, actualDir); err != nil {
			t.Fatalf("update golden dir %s: %v", goldenDir, err)
		}
	}

	goldenFiles := listRelativeFiles(t, goldenDir)
	actualFiles := listRelativeFiles(t, actualDir)
	if !reflect.DeepEqual(goldenFiles, actualFiles) {
		t.Fatalf("golden files mismatch\ngolden=%v\nactual=%v", goldenFiles, actualFiles)
	}

	for _, rel := range goldenFiles {
		assertGoldenFile(t, filepath.Join(goldenDir, rel), filepath.Join(actualDir, rel))
	}
}

func assertGoldenFile(t *testing.T, goldenPath, actualPath string) {
	t.Helper()

	goldenData, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden file %s: %v", goldenPath, err)
	}
	actualData, err := os.ReadFile(actualPath)
	if err != nil {
		t.Fatalf("read actual file %s: %v", actualPath, err)
	}

	if filepath.Ext(goldenPath) == ".json" {
		goldenData = normalizeJSON(t, goldenPath, goldenData)
		actualData = normalizeJSON(t, actualPath, actualData)
	}

	if !bytes.Equal(goldenData, actualData) {
		t.Fatalf("golden mismatch for %s\nwant:\n%s\ngot:\n%s", goldenPath, string(goldenData), string(actualData))
	}
}

func writeFixtureExportSettings(t *testing.T, caseDir, outputRoot string) string {
	t.Helper()

	settings := defaultFixtureExportSettings(outputRoot)
	path := filepath.Join(caseDir, "export_settings.json")
	if _, err := os.Stat(path); err == nil {
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			t.Fatalf("read fixture export settings %s: %v", path, readErr)
		}
		if err := json.Unmarshal(data, &settings); err != nil {
			t.Fatalf("unmarshal fixture export settings %s: %v", path, err)
		}
	} else if !os.IsNotExist(err) {
		t.Fatalf("stat fixture export settings %s: %v", path, err)
	}
	settings.OutputRoot = outputRoot

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		t.Fatalf("marshal fixture export settings: %v", err)
	}

	tempPath := filepath.Join(t.TempDir(), "export_settings.json")
	if err := os.WriteFile(tempPath, append(data, '\n'), 0o644); err != nil {
		t.Fatalf("write temp export settings: %v", err)
	}
	return tempPath
}

func defaultFixtureExportSettings(outputRoot string) config.ExportSettings {
	return config.ExportSettings{
		OutputRoot: outputRoot,
		ExportPaths: config.ExportPaths{
			GrimoireEffect: "generated/grimoire/effect",
			GrimoireDebug:  "generated/grimoire/give",
			ItemGive:       "generated/item/give",
			PassiveEffect:  "generated/passive/effect",
			PassiveGive:    "generated/passive/give",
			PassiveApply:   "generated/passive/apply",
			BowFlying:      "generated/bow/flying",
			BowGround:      "generated/bow/ground",
			Enemy:          "generated/enemy/spawn",
			EnemySkill:     "generated/enemy/skill",
			EnemyLoot:      "generated/enemy/loot",
		},
	}
}

func fixtureMinecraftLootRoot(caseDir string) string {
	return filepath.Join(caseDir, "minecraft_loot")
}

func syncGoldenDir(goldenDir, actualDir string) error {
	if err := os.RemoveAll(goldenDir); err != nil {
		return err
	}

	files, err := listRelativeFilesE(actualDir)
	if err != nil {
		return err
	}
	for _, rel := range files {
		src := filepath.Join(actualDir, rel)
		dst := filepath.Join(goldenDir, rel)
		data, err := os.ReadFile(src)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(dst, data, 0o644); err != nil {
			return err
		}
	}
	return nil
}

func listRelativeFiles(t *testing.T, root string) []string {
	t.Helper()

	files, err := listRelativeFilesE(root)
	if err != nil {
		t.Fatalf("list files %s: %v", root, err)
	}
	return files
}

func listRelativeFilesE(root string) ([]string, error) {
	if _, err := os.Stat(root); os.IsNotExist(err) {
		return []string{}, nil
	} else if err != nil {
		return nil, err
	}

	var files []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		files = append(files, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		return nil, err
	}
	if files == nil {
		return []string{}, nil
	}
	slices.Sort(files)
	return files, nil
}

func normalizeJSON(t *testing.T, path string, data []byte) []byte {
	t.Helper()

	var value any
	if err := json.Unmarshal(data, &value); err != nil {
		t.Fatalf("parse json %s: %v", path, err)
	}
	out, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatalf("marshal json %s: %v", path, err)
	}
	return append(out, '\n')
}

func compareStrings(a, b string) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}
