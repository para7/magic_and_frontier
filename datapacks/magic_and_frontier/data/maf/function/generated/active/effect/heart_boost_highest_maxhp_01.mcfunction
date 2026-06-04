function maf:common/target/highest_maxhp {range:10}
effect give @e[tag=maf_heal_target] minecraft:health_boost 300 1
execute at @e[tag=maf_heal_target] run particle minecraft:heart ~ ~1 ~ 0.3 0.3 0.3 1 6
function maf:common/target/clear
playsound minecraft:entity.player.levelup player @a ~ ~ ~ 2 1.5
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は ハートブースト を唱えた！"}]
