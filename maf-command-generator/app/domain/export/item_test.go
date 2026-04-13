package export

import (
	"path/filepath"
	"testing"
)

func TestItemExportFixtures(t *testing.T) {
	cases := discoverCases(t, filepath.Join("testdata", "item"))
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			master := loadFixtureMaster(t, tc.dir)
			artifacts, err := BuildItemArtifacts(master)
			if err != nil {
				t.Fatal(err)
			}

			actualDir := t.TempDir()
			if err := WriteItemArtifacts(filepath.Join(actualDir, "give"), artifacts); err != nil {
				t.Fatal(err)
			}

			assertGoldenDir(t, filepath.Join(tc.dir, "output"), actualDir)
		})
	}
}
