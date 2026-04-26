tag @e[type=#maf:enemymob,tag=maf_sight_lift_target] remove maf_sight_lift_target
tag @e[type=#maf:enemymob,tag=maf_sight_lift_candidate] remove maf_sight_lift_candidate
execute at @s anchored eyes rotated as @s positioned ^ ^ ^1 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^1 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^1.5 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^1.5 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^2 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^2 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^2.5 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^2.5 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^3 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^3 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^3.5 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^3.5 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^4 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^4 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^4.5 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^4.5 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^5 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^5 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^5.5 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^5.5 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^6 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^6 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^6.5 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^6.5 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^7 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^7 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^7.5 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^7.5 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^8 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^8 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^8.5 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^8.5 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^9 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^9 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^9.5 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^9.5 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^10 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^10 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^10.5 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^10.5 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^11 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^11 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^11.5 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^11.5 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^12 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^12 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^12.5 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^12.5 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^13 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^13 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^13.5 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^13.5 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^14 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^14 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^14.5 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^14.5 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^15 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^15 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^15.5 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^15.5 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^16 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^16 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^16.5 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^16.5 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^17 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^17 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^17.5 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^17.5 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^18 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^18 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^18.5 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^18.5 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^19 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^19 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^19.5 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^19.5 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^20 run tag @e[type=#maf:enemymob,distance=..1.3] add maf_sight_lift_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_lift_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^20 as @e[type=#maf:enemymob,distance=..1.3,sort=nearest,limit=1] run tag @s add maf_sight_lift_target
execute as @e[type=#maf:enemymob,tag=maf_sight_lift_target,distance=..25,sort=nearest,limit=1] run effect give @s minecraft:levitation 1 1
effect give @e[type=#maf:enemymob,tag=maf_sight_lift_candidate,distance=..25] minecraft:glowing 1 0 true
execute as @e[type=#maf:enemymob,tag=maf_sight_lift_target,distance=..25,sort=nearest,limit=1] at @s run particle minecraft:end_rod ~ ~0.9 ~ 0.2 0.4 0.2 0.01 18 force
tag @e[type=#maf:enemymob,tag=maf_sight_lift_candidate] remove maf_sight_lift_candidate
tag @e[type=#maf:enemymob,tag=maf_sight_lift_target] remove maf_sight_lift_target
playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 1.2 1.6
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は サイト・リフト を唱えた！"}]
