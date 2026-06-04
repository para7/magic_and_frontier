function maf:common/target/highest_armor {range:10}
execute as @e[tag=maf_heal_target] run function maf:buff/data/delay_heal/delay_heal_01/init
function maf:common/target/clear
playsound minecraft:entity.player.levelup player @a ~ ~ ~ 2 1.5
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は ディレイヒール を唱えた！"}]
