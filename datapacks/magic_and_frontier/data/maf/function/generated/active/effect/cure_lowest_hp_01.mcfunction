function maf:common/target/lowest_hp {range:10}
effect give @e[tag=maf_heal_target] minecraft:instant_health 1 0
execute at @e[tag=maf_heal_target] run particle minecraft:heart ~ ~1 ~ 0.3 0.3 0.3 1 4
function maf:common/target/clear
playsound minecraft:entity.player.levelup player @a ~ ~ ~ 2 2.0
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は キュア を唱えた！"}]
