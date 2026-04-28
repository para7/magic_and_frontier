tag @e[type=#maf:enemymob,tag=maf_sight_target] remove maf_sight_target
tag @e[type=#maf:enemymob,tag=maf_sight_candidate] remove maf_sight_candidate
execute at @s anchored eyes rotated as @s positioned ^ ^ ^0.5 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^0.5 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^0.5 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^0.5 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^0.5 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^1 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^1 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^1 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^1 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^1 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^1.5 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^1.5 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^1.5 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^1.5 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^1.5 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^2 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^2 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^2 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^2 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^2 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^2.5 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^2.5 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^2.5 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^2.5 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^2.5 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^3 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^3 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^3 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^3 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^3 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^3.5 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^3.5 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^3.5 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^3.5 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^3.5 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^4 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^4 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^4 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^4 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^4 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^4.5 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^4.5 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^4.5 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^4.5 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^4.5 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^5 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^5 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^5 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^5 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^5 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^5.5 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^5.5 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^5.5 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^5.5 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^5.5 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^6 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^6 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^6 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^6 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^6 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^6.5 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^6.5 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^6.5 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^6.5 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^6.5 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^7 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^7 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^7 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^7 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^7 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^7.5 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^7.5 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^7.5 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^7.5 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^7.5 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^8 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^8 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^8 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^8 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^8 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^8.5 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^8.5 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^8.5 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^8.5 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^8.5 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^9 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^9 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^9 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^9 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^9 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^9.5 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^9.5 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^9.5 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^9.5 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^9.5 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^10 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^10 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^10 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^10 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^10 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^10.5 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^10.5 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^10.5 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^10.5 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^10.5 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^11 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^11 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^11 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^11 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^11 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^11.5 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^11.5 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^11.5 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^11.5 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^11.5 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^12 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^12 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^12 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^12 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^12 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^12.5 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^12.5 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^12.5 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^12.5 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^12.5 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^13 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^13 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^13 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^13 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^13 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^13.5 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^13.5 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^13.5 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^13.5 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^13.5 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^14 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^14 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^14 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^14 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^14 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^14.5 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^14.5 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^14.5 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^14.5 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^14.5 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^15 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^15 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^15 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^15 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^15 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^15.5 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^15.5 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^15.5 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^15.5 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^15.5 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^16 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^16 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^16 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^16 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^16 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^16.5 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^16.5 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^16.5 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^16.5 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^16.5 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^17 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^17 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^17 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^17 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^17 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^17.5 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^17.5 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^17.5 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^17.5 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^17.5 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^18 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^18 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^18 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^18 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^18 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^18.5 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^18.5 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^18.5 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^18.5 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^18.5 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^19 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^19 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^19 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^19 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^19 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^19.5 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^19.5 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^19.5 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^19.5 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^19.5 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^20 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^20 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^20 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^20 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^20 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^20.5 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^20.5 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^20.5 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^20.5 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^20.5 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^21 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^21 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^21 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^21 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^21 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^21.5 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^21.5 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^21.5 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^21.5 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^21.5 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^22 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^22 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^22 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^22 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^22 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^22.5 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^22.5 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^22.5 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^22.5 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^22.5 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^23 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^23 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^23 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^23 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^23 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^23.5 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^23.5 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^23.5 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^23.5 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^23.5 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^24 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^24 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^24 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^24 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^24 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^24.5 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^24.5 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^24.5 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^24.5 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^24.5 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^25 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^25 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^25 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^25 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^25 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^25.5 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^25.5 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^25.5 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^25.5 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^25.5 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^26 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^26 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^26 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^26 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^26 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^26.5 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^26.5 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^26.5 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^26.5 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^26.5 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^27 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^27 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^27 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^27 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^27 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^27.5 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^27.5 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^27.5 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^27.5 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^27.5 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^28 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^28 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^28 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^28 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^28 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^28.5 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^28.5 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^28.5 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^28.5 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^28.5 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^29 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^29 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^29 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^29 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^29 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^29.5 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^29.5 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^29.5 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^29.5 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^29.5 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
execute at @s anchored eyes rotated as @s positioned ^ ^ ^30 if block ~ ~ ~ #maf:water run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^30 if block ~ ~ ~ minecraft:lava run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^30 unless block ~ ~ ~ #maf:sight_passable run return 0
execute at @s anchored eyes rotated as @s positioned ^ ^ ^30 as @e[type=#maf:enemymob,distance=..2.0] run tag @s add maf_sight_candidate
execute unless entity @e[type=#maf:enemymob,tag=maf_sight_target,limit=1] at @s anchored eyes rotated as @s positioned ^ ^ ^30 as @e[type=#maf:enemymob,distance=..2.0,sort=nearest,limit=1] run tag @s add maf_sight_target
