# 基礎値
scoreboard players set @s mafEquipMP 100

# === メインハンド ===
scoreboard players set @s tmp 0
execute if data entity @s SelectedItem.components."minecraft:custom_data".maf.maxmp run scoreboard players set @s tmp 1
execute if score @s tmp matches 1 store result score @s tmp2 run data get entity @s SelectedItem.components."minecraft:custom_data".maf.maxmp
execute if score @s tmp matches 1 run scoreboard players operation @s mafEquipMP += @s tmp2

# === オフハンド（Slot:-106b）===
scoreboard players set @s tmp 0
execute if data entity @s equipment.offhand.components."minecraft:custom_data".maf.maxmp run scoreboard players set @s tmp 1
execute if score @s tmp matches 1 store result score @s tmp2 run data get entity @s equipment.offhand.components."minecraft:custom_data".maf.maxmp
execute if score @s tmp matches 1 run scoreboard players operation @s mafEquipMP += @s tmp2

# === 頭部（Slot:103b）===
scoreboard players set @s tmp 0
execute if data entity @s equipment.head.components."minecraft:custom_data".maf.maxmp run scoreboard players set @s tmp 1
execute if score @s tmp matches 1 store result score @s tmp2 run data get entity @s equipment.head.components."minecraft:custom_data".maf.maxmp
execute if score @s tmp matches 1 run scoreboard players operation @s mafEquipMP += @s tmp2
execute if score @s tmp matches 0 if items entity @s armor.head #maf:leather_armor run scoreboard players add @s mafEquipMP 20
execute if score @s tmp matches 0 if items entity @s armor.head #maf:gold_armor run scoreboard players remove @s mafEquipMP 5
execute if score @s tmp matches 0 if items entity @s armor.head #maf:chain_armor run scoreboard players remove @s mafEquipMP 5
execute if score @s tmp matches 0 if items entity @s armor.head #maf:iron_armor run scoreboard players remove @s mafEquipMP 15
execute if score @s tmp matches 0 if items entity @s armor.head #maf:copper_armor run scoreboard players remove @s mafEquipMP 15
execute if score @s tmp matches 0 if items entity @s armor.head #maf:diamond_armor run scoreboard players remove @s mafEquipMP 20
execute if score @s tmp matches 0 if items entity @s armor.head #maf:netherite_armor run scoreboard players remove @s mafEquipMP 30

# === 胸部（Slot:102b）===
scoreboard players set @s tmp 0
execute if data entity @s equipment.chest.components."minecraft:custom_data".maf.maxmp run scoreboard players set @s tmp 1
execute if score @s tmp matches 1 store result score @s tmp2 run data get entity @s equipment.chest.components."minecraft:custom_data".maf.maxmp
execute if score @s tmp matches 1 run scoreboard players operation @s mafEquipMP += @s tmp2
execute if score @s tmp matches 0 if items entity @s armor.chest #maf:leather_armor run scoreboard players add @s mafEquipMP 20
execute if score @s tmp matches 0 if items entity @s armor.chest #maf:gold_armor run scoreboard players remove @s mafEquipMP 5
execute if score @s tmp matches 0 if items entity @s armor.chest #maf:chain_armor run scoreboard players remove @s mafEquipMP 5
execute if score @s tmp matches 0 if items entity @s armor.chest #maf:iron_armor run scoreboard players remove @s mafEquipMP 15
execute if score @s tmp matches 0 if items entity @s armor.chest #maf:copper_armor run scoreboard players remove @s mafEquipMP 15
execute if score @s tmp matches 0 if items entity @s armor.chest #maf:diamond_armor run scoreboard players remove @s mafEquipMP 20
execute if score @s tmp matches 0 if items entity @s armor.chest #maf:netherite_armor run scoreboard players remove @s mafEquipMP 30
execute if score @s tmp matches 0 if items entity @s armor.chest minecraft:elytra run scoreboard players remove @s mafEquipMP 250

# === 脚部（Slot:101b）===
scoreboard players set @s tmp 0
execute if data entity @s equipment.legs.components."minecraft:custom_data".maf.maxmp run scoreboard players set @s tmp 1
execute if score @s tmp matches 1 store result score @s tmp2 run data get entity @s equipment.legs.components."minecraft:custom_data".maf.maxmp
execute if score @s tmp matches 1 run scoreboard players operation @s mafEquipMP += @s tmp2
execute if score @s tmp matches 0 if items entity @s armor.legs #maf:leather_armor run scoreboard players add @s mafEquipMP 20
execute if score @s tmp matches 0 if items entity @s armor.legs #maf:gold_armor run scoreboard players remove @s mafEquipMP 5
execute if score @s tmp matches 0 if items entity @s armor.legs #maf:chain_armor run scoreboard players remove @s mafEquipMP 5
execute if score @s tmp matches 0 if items entity @s armor.legs #maf:iron_armor run scoreboard players remove @s mafEquipMP 15
execute if score @s tmp matches 0 if items entity @s armor.legs #maf:copper_armor run scoreboard players remove @s mafEquipMP 15
execute if score @s tmp matches 0 if items entity @s armor.legs #maf:diamond_armor run scoreboard players remove @s mafEquipMP 20
execute if score @s tmp matches 0 if items entity @s armor.legs #maf:netherite_armor run scoreboard players remove @s mafEquipMP 30

# === 足部（Slot:100b）===
scoreboard players set @s tmp 0
execute if data entity @s equipment.feet.components."minecraft:custom_data".maf.maxmp run scoreboard players set @s tmp 1
execute if score @s tmp matches 1 store result score @s tmp2 run data get entity @s equipment.feet.components."minecraft:custom_data".maf.maxmp
execute if score @s tmp matches 1 run scoreboard players operation @s mafEquipMP += @s tmp2
execute if score @s tmp matches 0 if items entity @s armor.feet #maf:leather_armor run scoreboard players add @s mafEquipMP 20
execute if score @s tmp matches 0 if items entity @s armor.feet #maf:gold_armor run scoreboard players remove @s mafEquipMP 5
execute if score @s tmp matches 0 if items entity @s armor.feet #maf:chain_armor run scoreboard players remove @s mafEquipMP 5
execute if score @s tmp matches 0 if items entity @s armor.feet #maf:iron_armor run scoreboard players remove @s mafEquipMP 15
execute if score @s tmp matches 0 if items entity @s armor.feet #maf:copper_armor run scoreboard players remove @s mafEquipMP 15
execute if score @s tmp matches 0 if items entity @s armor.feet #maf:diamond_armor run scoreboard players remove @s mafEquipMP 20
execute if score @s tmp matches 0 if items entity @s armor.feet #maf:netherite_armor run scoreboard players remove @s mafEquipMP 30

# 最低値を0にクランプ
execute if score @s mafEquipMP matches ..-1 run scoreboard players set @s mafEquipMP 0
