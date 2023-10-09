tellraw @s [{"text":"初期化処理です"}]

scoreboard players set @s p7_playerID 0
scoreboard players set @s p7_move 0
scoreboard players set @s p7_isMovedY 0
scoreboard players set @s p7_posXpre 0
scoreboard players set @s p7_posYpre 0
scoreboard players set @s p7_posZpre 0
scoreboard players set @s p7_posX 0
scoreboard players set @s p7_posY 0
scoreboard players set @s p7_posZ 0
scoreboard players set @s p7_soul 0
scoreboard players set @s p7_soulTick 0
# scoreboard players set @s const0 0
scoreboard players set @s tmp 0
scoreboard players set @s tmp2 0

scoreboard players set @s mafLogin 1

function maf:magic/constructor