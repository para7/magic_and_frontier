tellraw @s [{"text":"詠唱が中断されました"}]
# TODO: ログアウト時に詠唱を中断させる
scoreboard players set @s p7_castCost 0
scoreboard players set @s p7_castTime -1
scoreboard players set @s p7_castTimeMax 0
scoreboard players set @s p7_castID 0