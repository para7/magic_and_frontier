package treasures

import (
	"strings"

	"maf-command-editor/app/internal/domain/common"
)

func IsSupportedTablePath(tablePath string) bool {
	path := common.NormalizeText(tablePath)
	return strings.HasPrefix(path, "minecraft:chests/")
}
