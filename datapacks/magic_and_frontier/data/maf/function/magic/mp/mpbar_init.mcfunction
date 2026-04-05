# マクロ関数: bossbar初期化 (load時1回のみ)  引数: $(id) = プレイヤーID (1..20)
$bossbar add mpbar$(id) "MP"
$bossbar set minecraft:mpbar$(id) style progress
$bossbar set minecraft:mpbar$(id) color green
