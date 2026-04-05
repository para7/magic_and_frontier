# マクロ関数: プレイヤー単体のMPbar更新  引数: $(mpbar_id)  実行者: @s = 対象プレイヤー
$bossbar set minecraft:mpbar$(mpbar_id) players @s
$execute if score @s mafCastTime matches ..0 run bossbar set minecraft:mpbar$(mpbar_id) color green
$execute if score @s mafCastTime matches ..0 run bossbar set minecraft:mpbar$(mpbar_id) style progress
$execute if score @s mafCastTime matches ..0 run bossbar set minecraft:mpbar$(mpbar_id) name [{"text":"MP "},{"score":{"name":"@s","objective":"mafMP"}},{"text": " / "}, {"score":{"name":"@s","objective":"mafMaxMP"}}]
$execute if score @s mafCastTime matches ..0 run execute store result bossbar minecraft:mpbar$(mpbar_id) max run scoreboard players get @s mafMaxMP
$execute if score @s mafCastTime matches ..0 run execute store result bossbar minecraft:mpbar$(mpbar_id) value run scoreboard players get @s mafMP
$execute if score @s mafCastTime matches 1.. run bossbar set minecraft:mpbar$(mpbar_id) color yellow
$execute if score @s mafCastTime matches 1.. run bossbar set minecraft:mpbar$(mpbar_id) style notched_20
$execute if score @s mafCastTime matches 1.. run bossbar set minecraft:mpbar$(mpbar_id) name [{"text":"Casting... "},{"storage":"p7:mpbar","nbt":"bar$(mpbar_id).title"}]
$execute if score @s mafCastTime matches 1.. run execute store result bossbar minecraft:mpbar$(mpbar_id) max run scoreboard players get @s mafCastTimeMax
$execute if score @s mafCastTime matches 1.. run execute store result bossbar minecraft:mpbar$(mpbar_id) value run scoreboard players get @s mafCastTime
