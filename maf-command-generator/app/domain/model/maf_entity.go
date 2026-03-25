package model

import (
	master "maf_command_editor/app/domain/master"
)

type MafEntity[T any] interface {
	ValidateJSON(data any, mas master.DBMaster) (T, error)

	// 主に web 画面からの操作用、メモリで保持してる配列にデータを追記する。リレーション関係を確認するため validate も行う
	Create(data T, mas master.DBMaster) error
	Update(data T, mas master.DBMaster) error
	Delete(id string, mas master.DBMaster) error
	Save() error
	Load() error

	// DBMaster との連携用
	// ListAll() []T
	ValidateAll(mas master.DBMaster) []error
	Find() T
}
