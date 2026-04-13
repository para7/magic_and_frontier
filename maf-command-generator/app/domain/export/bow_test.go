package export

import (
	"path/filepath"
	"testing"
)

func TestBowExportFixtures(t *testing.T) {
	cases := discoverCases(t, filepath.Join("testdata", "bow"))
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			master := loadFixtureMaster(t, tc.dir)
			effects, hits, flyings, grounds, err := BuildBowArtifacts(master)
			if err != nil {
				t.Fatal(err)
			}

			actualDir := t.TempDir()
			if err := WriteBowArtifacts(
				filepath.Join(actualDir, "effect"),
				filepath.Join(actualDir, "bow"),
				filepath.Join(actualDir, "flying"),
				filepath.Join(actualDir, "ground"),
				master.ListBows(),
				effects,
				hits,
				flyings,
				grounds,
			); err != nil {
				t.Fatal(err)
			}

			assertGoldenDir(t, filepath.Join(tc.dir, "output"), actualDir)
		})
	}
}
