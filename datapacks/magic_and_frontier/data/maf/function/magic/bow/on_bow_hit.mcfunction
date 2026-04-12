# フラグをリセット
scoreboard players set @s mafBowHit 0

execute store result storage maf:tmp bow_player_id int 1 run scoreboard players get @s mafPlayerID
function maf:magic/bow/process_hit_arrows with storage maf:tmp
data remove storage maf:tmp bow_player_id
