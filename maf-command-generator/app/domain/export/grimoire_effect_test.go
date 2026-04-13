package export

import (
	"path/filepath"
	"testing"
)

func TestGrimoireExportFixtures(t *testing.T) {
	cases := discoverCases(t, filepath.Join("testdata", "grimoire"))
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			master := loadFixtureMaster(t, tc.dir)
			effects := BuildGrimoireArtifacts(master)

			actualDir := t.TempDir()
			if err := WriteGrimoireArtifacts(filepath.Join(actualDir, "effect"), effects); err != nil {
				t.Fatal(err)
			}
			if err := WriteGrimoireDebugArtifacts(filepath.Join(actualDir, "give"), effects); err != nil {
				t.Fatal(err)
			}

			assertGoldenDir(t, filepath.Join(tc.dir, "output"), actualDir)
		})
	}
}
