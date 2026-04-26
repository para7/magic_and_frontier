tag @e[type=#maf:enemymob,tag=maf_sight_lift_target] remove maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^1 as @e[type=#maf:enemymob,distance=..1.2] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^2 as @e[type=#maf:enemymob,distance=..1.2] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^3 as @e[type=#maf:enemymob,distance=..1.2] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^4 as @e[type=#maf:enemymob,distance=..1.2] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^5 as @e[type=#maf:enemymob,distance=..1.2] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^6 as @e[type=#maf:enemymob,distance=..1.2] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^7 as @e[type=#maf:enemymob,distance=..1.2] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^8 as @e[type=#maf:enemymob,distance=..1.2] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^9 as @e[type=#maf:enemymob,distance=..1.2] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^10 as @e[type=#maf:enemymob,distance=..1.2] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^11 as @e[type=#maf:enemymob,distance=..1.2] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^12 as @e[type=#maf:enemymob,distance=..1.2] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^13 as @e[type=#maf:enemymob,distance=..1.2] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^14 as @e[type=#maf:enemymob,distance=..1.2] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^15 as @e[type=#maf:enemymob,distance=..1.2] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^16 as @e[type=#maf:enemymob,distance=..1.2] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^17 as @e[type=#maf:enemymob,distance=..1.2] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^18 as @e[type=#maf:enemymob,distance=..1.2] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^19 as @e[type=#maf:enemymob,distance=..1.2] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^20 as @e[type=#maf:enemymob,distance=..1.2] run tag @s add maf_sight_lift_target
execute at @s as @e[type=#maf:enemymob,tag=maf_sight_lift_target,distance=..25,sort=nearest,limit=1] run effect give @s minecraft:levitation 3 1
execute at @s as @e[type=#maf:enemymob,tag=maf_sight_lift_target,distance=..25,sort=nearest,limit=1] at @s run particle minecraft:end_rod ~ ~0.9 ~ 0.2 0.4 0.2 0.01 18 force
tag @e[type=#maf:enemymob,tag=maf_sight_lift_target] remove maf_sight_lift_target
playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 1.2 1.6
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は サイト・リフト を唱えた！"}]
