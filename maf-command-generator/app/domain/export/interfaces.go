package export

import (
	grimoireModel "maf_command_editor/app/domain/model/grimoire"
)

type DBMaster interface {
	GetGrimoireByID(id string) (grimoireModel.Grimoire, bool)
	ListGrimoires() []grimoireModel.Grimoire
}
