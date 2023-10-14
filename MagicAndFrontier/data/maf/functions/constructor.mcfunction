tellraw @s [{"text":"初期化処理です"}]

scoreboard players set @s mafPlayerID 0
scoreboard players set @s mafMoved 0
scoreboard players set @s mafIsMovedY 0
# scoreboard players set @s p7_posXpre 0
scoreboard players set @s mafPosYpre 0
# scoreboard players set @s p7_posZpre 0
# scoreboard players set @s p7_posX 0
scoreboard players set @s mafPosY 0
# scoreboard players set @s p7_posZ 0
scoreboard players set @s mafSoul 0
scoreboard players set @s mafSoulTick 0
# scoreboard players set @s const0 0
scoreboard players set @s tmp 0
scoreboard players set @s tmp2 0

scoreboard players set @s mafLogin 1
scoreboard players set @s mafCastTime -1

function maf:magic/constructor