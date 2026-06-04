function maf:common/buff/set {buff_id:"delay_heal_01",buff_category:"delay_heal",tick:300}
tag @s add maf_has_buff
execute at @s run particle minecraft:heart ~ ~1.0 ~ 0.3 0.3 0.3 0.5 3 force
playsound minecraft:entity.player.attack.sweep master @a[distance=..24] ~ ~ ~ 1.2 1.5
