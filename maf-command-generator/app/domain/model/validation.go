package model

type ValidationError struct {
	Entity string // "grimoire", "item" など
	ID     string // エンティティのID
	Field  string // エラーのあったフィールド名（jsonタグ名）
	Tag    string // validatorのタグ名
	Param  string // validatorのタグパラメータ
}
