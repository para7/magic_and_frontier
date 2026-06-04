function maf:common/buff/set {buff_id:"trigger_heal_01",buff_category:"trigger_heal",tick:600}
tag @s add maf_has_buff
execute at @s run particle minecraft:heart ~ ~1.0 ~ 0.3 0.3 0.3 0.5 3 force
playsound minecraft:entity.player.attack.sweep master @a[distance=..24] ~ ~ ~ 1.2 1.5
