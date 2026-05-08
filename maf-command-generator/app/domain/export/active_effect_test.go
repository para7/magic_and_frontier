package export

import (
	"path/filepath"
	"testing"
)

func TestActiveExportFixtures(t *testing.T) {
	cases := discoverCases(t, filepath.Join("testdata", "active"))
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			master := loadFixtureMaster(t, tc.dir)
			effects := BuildActiveArtifacts(master)

			actualDir := t.TempDir()
			if err := WriteActiveArtifacts(filepath.Join(actualDir, "effect"), effects); err != nil {
				t.Fatal(err)
			}
			if err := WriteActiveSpellArtifacts(filepath.Join(actualDir, "spell"), effects); err != nil {
				t.Fatal(err)
			}
			if err := WriteActiveDebugArtifacts(filepath.Join(actualDir, "give"), effects); err != nil {
				t.Fatal(err)
			}

			assertGoldenDir(t, filepath.Join(tc.dir, "output"), actualDir)
		})
	}
}
