package export

import (
	"path/filepath"
	"testing"
)

func TestPassiveExportFixtures(t *testing.T) {
	cases := discoverCases(t, filepath.Join("testdata", "passive"))
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			master := loadFixtureMaster(t, tc.dir)
			effects, bows, grimoires, err := BuildPassiveArtifacts(master)
			if err != nil {
				t.Fatal(err)
			}

			actualDir := t.TempDir()
			if err := WritePassiveArtifacts(
				filepath.Join(actualDir, "effect"),
				filepath.Join(actualDir, "bow"),
				filepath.Join(actualDir, "give"),
				filepath.Join(actualDir, "apply"),
				effects,
				bows,
				grimoires,
			); err != nil {
				t.Fatal(err)
			}

			assertGoldenDir(t, filepath.Join(tc.dir, "output"), actualDir)
		})
	}
}
