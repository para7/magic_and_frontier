# effect give @e[type=#maf:undead,sort=nearest,limit=1] minecraft:instant_health 1 1
# effect give @e[type=#maf:undead,sort=nearest,limit=1] minecraft:slowness 14 40
# effect give @e[type=#maf:undead,sort=nearest,limit=1] minecraft:glowing 14 0
# effect give @e[type=#maf:undead,sort=nearest,limit=1] minecraft:weakness 14 40
execute at @e[distance=0.1..18,type=#maf:friendmob] run particle minecraft:heart ~ ~1 ~ 0.5 0.5 0.5 1 6
effect give @e[distance=0.1..18,type=#maf:friendmob] minecraft:regeneration 4 14
effect give @s minecraft:instant_damage 1 1
particle minecraft:heart ~ ~ ~ 0.5 0.5 0.5 1 5
playsound minecraft:entity.player.levelup player @a ~ ~ ~ 1 2.0
playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 2
playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 0.5
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は サクリファイス を唱えた！"}]
