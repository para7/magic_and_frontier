tellraw @s [{"text":"set_magic"}]

# 0になった時発動とするので、マイナスで初期化してバグ対策
scoreboard players set @s p7_castTime -1

# 一時変数からデータをロード
execute store result score @s p7_castCost run data get storage p7:maf magictmp.cost
execute store result score @s p7_castTime run data get storage p7:maf magictmp.cast
execute store result score @s p7_castID run data get storage p7:maf magictmp.id

execute if score @s p7_castCost > @s p7_MP run scoreboard players set @s p7_castTime -1
execute if score @s p7_castCost > @s p7_MP run tellraw @s [{"text":"MPが足りません！"}, {"score":{"name":"@s","objective":"p7_MP"}},{"text":" / "},{"score":{"name":"@s","objective":"p7_castCost"}}]
execute if score @s p7_castCost > @s p7_MP run playsound minecraft:block.dispenser.fail master @s ~ ~ ~ 1.0 1.1

# 発動条件のある魔法はここで判定をする？
# 詠唱後でいいかも