function maf:common/target/most_hp_missing {range:10}
execute as @e[tag=maf_heal_target] run function maf:buff/data/trigger_heal/trigger_heal_01/init
function maf:common/target/clear
playsound minecraft:entity.player.levelup player @a ~ ~ ~ 2 1.5
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は トリガーヒール を唱えた！"}]
