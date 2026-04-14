package export

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPassiveExportFixtures(t *testing.T) {
	cases := discoverCases(t, filepath.Join("testdata", "passive"))
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			master := loadFixtureMaster(t, tc.dir)
			effects, grimoires, err := BuildPassiveArtifacts(master)
			if err != nil {
				t.Fatal(err)
			}

			actualDir := t.TempDir()
			if err := WritePassiveArtifacts(
				filepath.Join(actualDir, "effect"),
				filepath.Join(actualDir, "give"),
				filepath.Join(actualDir, "apply"),
				effects,
				grimoires,
			); err != nil {
				t.Fatal(err)
			}

			assertGoldenDir(t, filepath.Join(tc.dir, "output"), actualDir)
		})
	}
}

func TestWritePassiveArtifactsRemovesStaleSlotFunctions(t *testing.T) {
	root := t.TempDir()
	giveDir := filepath.Join(root, "give")
	applyDir := filepath.Join(root, "apply")

	if err := writeFunctionFile(filepath.Join(giveDir, "stale_slot1.mcfunction"), ""); err != nil {
		t.Fatal(err)
	}
	if err := writeFunctionFile(filepath.Join(applyDir, "stale_slot1.mcfunction"), ""); err != nil {
		t.Fatal(err)
	}

	if err := WritePassiveArtifacts(
		filepath.Join(root, "effect"),
		giveDir,
		applyDir,
		[]PassiveEffectFunction{{ID: "passive_1", Body: "say effect"}},
		nil,
	); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(giveDir, "stale_slot1.mcfunction")); !os.IsNotExist(err) {
		t.Fatalf("expected stale give function to be removed, err=%v", err)
	}
	if _, err := os.Stat(filepath.Join(applyDir, "stale_slot1.mcfunction")); !os.IsNotExist(err) {
		t.Fatalf("expected stale apply function to be removed, err=%v", err)
	}
}
