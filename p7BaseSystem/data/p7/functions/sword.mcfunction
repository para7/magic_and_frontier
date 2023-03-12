
scoreboard players set @a p7UseSword 0
execute if entity @a[scores={p7dealt=1..},nbt={SelectedItem:{id:"minecraft:wooden_sword"}}] run scoreboard players set @a p7UseSword 1
execute if entity @a[scores={p7dealt=1..},nbt={SelectedItem:{id:"minecraft:iron_sword"}}] run scoreboard players set @a p7UseSword 1
execute if entity @a[scores={p7dealt=1..},nbt={SelectedItem:{id:"minecraft:golden_sword"}}] run scoreboard players set @a p7UseSword 1
execute if entity @a[scores={p7dealt=1..},nbt={SelectedItem:{id:"minecraft:diamond_sword"}}] run scoreboard players set @a p7UseSword 1
execute if entity @a[scores={p7dealt=1..},nbt={SelectedItem:{id:"minecraft:netherite_sword"}}] run scoreboard players set @a p7UseSword 1

# execute if entity @a[scores={p7UseSword=1..}] run tell @a sword
scoreboard players set @a p7dealt 0