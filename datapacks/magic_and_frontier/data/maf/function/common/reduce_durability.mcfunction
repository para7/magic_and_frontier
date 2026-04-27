# マクロ関数: メインハンド装備の耐久を減らす  引数: $(amount)  実行者: @s = 対象プレイヤー
execute unless data entity @s SelectedItem run return 0

# 一時 armor_stand を演算領域として使い、item stack 全体を戻す
execute at @s run summon minecraft:armor_stand ~ ~ ~ {Invisible:1b,Marker:1b,NoGravity:1b,Invulnerable:1b,Silent:1b,Tags:["maf_reduce_durability_tmp"]}
execute at @s run item replace entity @e[type=minecraft:armor_stand,tag=maf_reduce_durability_tmp,distance=..0.5,sort=nearest,limit=1] weapon.mainhand from entity @s weapon.mainhand
execute at @s run item replace entity @e[type=minecraft:armor_stand,tag=maf_reduce_durability_tmp,distance=..0.5,sort=nearest,limit=1] weapon.offhand from entity @s weapon.mainhand

# 現在ダメージ (未破損なら 0)
scoreboard players set #reduce_durability_current tmp 0
execute store result score #reduce_durability_current tmp run data get entity @s SelectedItem.components."minecraft:damage"

# 最大耐久を実測するため full damage を適用して読み取る
execute at @s run item modify entity @e[type=minecraft:armor_stand,tag=maf_reduce_durability_tmp,distance=..0.5,sort=nearest,limit=1] weapon.mainhand maf:common/reduce_durability_set_full
scoreboard players set #reduce_durability_max tmp 0
execute at @s store result score #reduce_durability_max tmp run data get entity @e[type=minecraft:armor_stand,tag=maf_reduce_durability_tmp,distance=..0.5,sort=nearest,limit=1] equipment.mainhand.components."minecraft:damage"
execute unless score #reduce_durability_max tmp matches 1.. at @s run kill @e[type=minecraft:armor_stand,tag=maf_reduce_durability_tmp,distance=..0.5,sort=nearest,limit=1]
execute unless score #reduce_durability_max tmp matches 1.. run return 0

# new_damage = current + amount, max以上なら破壊
$scoreboard players add #reduce_durability_current tmp $(amount)
execute if score #reduce_durability_current tmp >= #reduce_durability_max tmp run playsound minecraft:entity.item.break player @s ~ ~ ~ 1.2 1.0
execute if score #reduce_durability_current tmp >= #reduce_durability_max tmp run item replace entity @s weapon.mainhand with air
execute if score #reduce_durability_current tmp >= #reduce_durability_max tmp at @s run kill @e[type=minecraft:armor_stand,tag=maf_reduce_durability_tmp,distance=..0.5,sort=nearest,limit=1]
execute if score #reduce_durability_current tmp >= #reduce_durability_max tmp run return 0

# 退避していた item stack の damage を整数で更新して mainhand を丸ごと戻す
execute at @s store result entity @e[type=minecraft:armor_stand,tag=maf_reduce_durability_tmp,distance=..0.5,sort=nearest,limit=1] equipment.offhand.components."minecraft:damage" int 1 run scoreboard players get #reduce_durability_current tmp
execute at @s run item replace entity @s weapon.mainhand from entity @e[type=minecraft:armor_stand,tag=maf_reduce_durability_tmp,distance=..0.5,sort=nearest,limit=1] weapon.offhand
execute at @s run kill @e[type=minecraft:armor_stand,tag=maf_reduce_durability_tmp,distance=..0.5,sort=nearest,limit=1]
