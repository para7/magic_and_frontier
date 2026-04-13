package export

import (
	"path/filepath"
	"testing"
)

func TestEnemySkillExportFixtures(t *testing.T) {
	cases := discoverCases(t, filepath.Join("testdata", "enemy_skill"))
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			master := loadFixtureMaster(t, tc.dir)
			artifacts := BuildEnemySkillArtifacts(master, "generated/enemy/skill")

			actualDir := t.TempDir()
			if err := WriteEnemySkillArtifacts(filepath.Join(actualDir, "skill"), artifacts); err != nil {
				t.Fatal(err)
			}

			assertGoldenDir(t, filepath.Join(tc.dir, "output"), actualDir)
		})
	}
}

func TestEnemyExportFixtures(t *testing.T) {
	cases := discoverCases(t, filepath.Join("testdata", "enemy"))
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			master := loadFixtureMaster(t, tc.dir)
			artifacts, err := BuildEnemyArtifacts(master, "generated/enemy/loot", fixtureMinecraftLootRoot(tc.dir))
			if err != nil {
				t.Fatal(err)
			}

			actualDir := t.TempDir()
			if err := WriteEnemyArtifacts(filepath.Join(actualDir, "spawn"), filepath.Join(actualDir, "loot"), artifacts); err != nil {
				t.Fatal(err)
			}

			assertGoldenDir(t, filepath.Join(tc.dir, "output"), actualDir)
		})
	}
}
