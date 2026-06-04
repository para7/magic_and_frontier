function maf:common/target/lowest_maxhp {range:10}
effect give @e[tag=maf_heal_target] minecraft:absorption 210 2
execute at @e[tag=maf_heal_target] run particle minecraft:end_rod ~ ~1 ~ 0.3 0.5 0.3 0.05 8
function maf:common/target/clear
playsound minecraft:entity.player.levelup player @a ~ ~ ~ 2 1.5
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は バリア を唱えた！"}]
