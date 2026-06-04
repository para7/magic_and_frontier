scoreboard players set @s tmp 0
execute if entity @s[nbt={active_effects:[{id:"minecraft:poison"}]}] run scoreboard players add @s tmp 1
execute if entity @s[nbt={active_effects:[{id:"minecraft:poison"}]}] run effect clear @s minecraft:poison
execute if entity @s[nbt={active_effects:[{id:"minecraft:slowness"}]}] run scoreboard players add @s tmp 1
execute if entity @s[nbt={active_effects:[{id:"minecraft:slowness"}]}] run effect clear @s minecraft:slowness
execute if entity @s[nbt={active_effects:[{id:"minecraft:mining_fatigue"}]}] run scoreboard players add @s tmp 1
execute if entity @s[nbt={active_effects:[{id:"minecraft:mining_fatigue"}]}] run effect clear @s minecraft:mining_fatigue
execute if entity @s[nbt={active_effects:[{id:"minecraft:blindness"}]}] run scoreboard players add @s tmp 1
execute if entity @s[nbt={active_effects:[{id:"minecraft:blindness"}]}] run effect clear @s minecraft:blindness
execute if entity @s[nbt={active_effects:[{id:"minecraft:wither"}]}] run scoreboard players add @s tmp 1
execute if entity @s[nbt={active_effects:[{id:"minecraft:wither"}]}] run effect clear @s minecraft:wither
execute if entity @s[nbt={active_effects:[{id:"minecraft:hunger"}]}] run scoreboard players add @s tmp 1
execute if entity @s[nbt={active_effects:[{id:"minecraft:hunger"}]}] run effect clear @s minecraft:hunger
execute if entity @s[nbt={active_effects:[{id:"minecraft:levitation"}]}] run scoreboard players add @s tmp 1
execute if entity @s[nbt={active_effects:[{id:"minecraft:levitation"}]}] run effect clear @s minecraft:levitation
execute if score @s tmp matches 1.. run effect give @s minecraft:instant_health 1 0
execute if score @s tmp matches 3.. run effect give @s minecraft:instant_health 1 1
execute if score @s tmp matches 5.. run effect give @s minecraft:instant_health 1 2
particle minecraft:heart ~ ~1 ~ 0.5 0.5 0.5 1 8
playsound minecraft:entity.player.levelup player @a ~ ~ ~ 2 1.0
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は カースコンバート を唱えた！"}]
