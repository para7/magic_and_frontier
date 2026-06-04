function maf:common/buff/set {buff_id:"resist_curse_01",buff_category:"resist_curse",tick:600}
tag @s add maf_has_buff
execute at @s run particle minecraft:enchanted_hit ~ ~1.0 ~ 0.35 0.45 0.35 0.05 12 force
playsound minecraft:entity.player.attack.sweep master @a[distance=..24] ~ ~ ~ 1.2 1.5
