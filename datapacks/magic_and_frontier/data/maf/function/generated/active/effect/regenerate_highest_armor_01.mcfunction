function maf:common/target/highest_armor {range:10}
effect give @e[tag=maf_heal_target] minecraft:regeneration 30 1
execute at @e[tag=maf_heal_target] run particle minecraft:heart ~ ~1 ~ 0.3 0.3 0.3 0.5 6
function maf:common/target/clear
playsound minecraft:entity.player.levelup player @a ~ ~ ~ 2 2.0
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は リジェネレート を唱えた！"}]
