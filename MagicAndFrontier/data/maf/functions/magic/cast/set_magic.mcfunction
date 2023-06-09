tellraw @s [{"text":"set_magic"}]

# 0になった時発動とするので、マイナスで初期化してバグ対策
scoreboard players set @s p7_castTime -1

# 一時変数からデータをロード
execute store result score @s p7_castCost run data get storage p7:maf magictmp.cost
execute store result score @s p7_castTime run data get storage p7:maf magictmp.cast
execute store result score @s p7_castTimeMax run data get storage p7:maf magictmp.cast
execute store result score @s p7_castID run data get storage p7:maf magictmp.id

execute if score @s p7_castCost > @s p7_MP run scoreboard players set @s p7_castTime -1
# execute if score @s p7_castCost > @s p7_MP run tellraw @s [{"text":"MPが足りません！"}, {"score":{"name":"@s","objective":"p7_MP"}},{"text":" / "},{"score":{"name":"@s","objective":"p7_castCost"}}]
execute if score @s p7_castCost > @s p7_MP run tellraw @s [{"text":"MPが足りません！ 消費MP: "},{"score":{"name":"@s","objective":"p7_castCost"}}]
execute if score @s p7_castCost > @s p7_MP run playsound minecraft:block.dispenser.fail master @s ~ ~ ~ 1.0 1.1

# 詠唱名を保存する
execute as @a[scores={p7_playerID=1}] run data modify storage p7:mpbar bar1.title set from storage p7:maf magictmp.title
execute as @a[scores={p7_playerID=2}] run data modify storage p7:mpbar bar2.title set from storage p7:maf magictmp.title
execute as @a[scores={p7_playerID=3}] run data modify storage p7:mpbar bar3.title set from storage p7:maf magictmp.title
execute as @a[scores={p7_playerID=4}] run data modify storage p7:mpbar bar4.title set from storage p7:maf magictmp.title
execute as @a[scores={p7_playerID=5}] run data modify storage p7:mpbar bar5.title set from storage p7:maf magictmp.title
execute as @a[scores={p7_playerID=6}] run data modify storage p7:mpbar bar6.title set from storage p7:maf magictmp.title
execute as @a[scores={p7_playerID=7}] run data modify storage p7:mpbar bar7.title set from storage p7:maf magictmp.title
execute as @a[scores={p7_playerID=8}] run data modify storage p7:mpbar bar8.title set from storage p7:maf magictmp.title
execute as @a[scores={p7_playerID=9}] run data modify storage p7:mpbar bar9.title set from storage p7:maf magictmp.title
execute as @a[scores={p7_playerID=10}] run data modify storage p7:mpbar bar10.title set from storage p7:maf magictmp.title
execute as @a[scores={p7_playerID=11}] run data modify storage p7:mpbar bar11.title set from storage p7:maf magictmp.title
execute as @a[scores={p7_playerID=12}] run data modify storage p7:mpbar bar12.title set from storage p7:maf magictmp.title
execute as @a[scores={p7_playerID=13}] run data modify storage p7:mpbar bar13.title set from storage p7:maf magictmp.title
execute as @a[scores={p7_playerID=14}] run data modify storage p7:mpbar bar14.title set from storage p7:maf magictmp.title
execute as @a[scores={p7_playerID=15}] run data modify storage p7:mpbar bar15.title set from storage p7:maf magictmp.title
execute as @a[scores={p7_playerID=16}] run data modify storage p7:mpbar bar16.title set from storage p7:maf magictmp.title
execute as @a[scores={p7_playerID=17}] run data modify storage p7:mpbar bar17.title set from storage p7:maf magictmp.title
execute as @a[scores={p7_playerID=18}] run data modify storage p7:mpbar bar18.title set from storage p7:maf magictmp.title
execute as @a[scores={p7_playerID=19}] run data modify storage p7:mpbar bar19.title set from storage p7:maf magictmp.title
execute as @a[scores={p7_playerID=20}] run data modify storage p7:mpbar bar20.title set from storage p7:maf magictmp.title


# 発動条件のある魔法はここで判定をする？
# 発動時でいいかも　「力が足りなかった！」的な
