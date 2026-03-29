package cli

import (
	"fmt"
	"maf_command_editor/app/domain/export"
	master "maf_command_editor/app/domain/master"
	config "maf_command_editor/app/files"
)

func Export(dmas *master.DBMaster, cfg config.MafConfig) int {
	// export を実行していいかを確認する
	result := Validate(dmas)

	if result != 0 {
		fmt.Print("データに問題があります。エクスポートを中断します。\n")
		return result
	}

	err := export.ExportDatapack(dmas, cfg)

	if err != nil {
		fmt.Printf("エクスポート中にエラーが発生しました: %v\n", err)
		return 1
	}

	fmt.Print("エクスポートが完了しました。\n")
	return 0
}
