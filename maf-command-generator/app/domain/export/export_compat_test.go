package export

import (
	"path/filepath"
	"testing"

	config "maf_command_editor/app/files"
)

func TestExportDatapackFixtures(t *testing.T) {
	cases := discoverCases(t, filepath.Join("testdata", "export"))
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			master := loadFixtureMaster(t, tc.dir)

			outputRoot := filepath.Join(t.TempDir(), "out")
			cfg := config.LoadConfig()
			cfg.ExportSettingsPath = writeFixtureExportSettings(t, tc.dir, outputRoot)
			cfg.MinecraftLootTableRoot = fixtureMinecraftLootRoot(tc.dir)
			cfg.LootTableSourceRoot = filepath.Join(tc.dir, "input", "loot_table")

			if err := ExportDatapack(master, cfg); err != nil {
				t.Fatal(err)
			}

			assertGoldenDir(t, filepath.Join(tc.dir, "output"), filepath.Join(outputRoot, "data", "maf"))
		})
	}
}
