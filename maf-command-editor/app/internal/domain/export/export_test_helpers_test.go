package export

import (
	"encoding/json"
	"os"
	"testing"
)

func writeExportSettingsFile(t *testing.T, path string, settings ExportSettings) {
	t.Helper()
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}
