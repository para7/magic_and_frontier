package treasures

import (
	"strings"

	"tools2/app/internal/domain/common"
)

func IsSupportedTablePath(tablePath string) bool {
	path := common.NormalizeText(tablePath)
	return strings.HasPrefix(path, "minecraft:chests/")
}
