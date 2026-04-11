# フラグをリセット
scoreboard players set @s mafBowHit 0

# タグ付き矢が近くになければ終了
execute unless entity @e[type=arrow,tag=maf_passive_arrow,limit=1,sort=nearest,distance=..10] run return 0

# 矢のパッシブ情報を退避（kill 前に取得）
data modify storage maf:tmp arrow_passive set from entity @e[type=arrow,tag=maf_passive_arrow,sort=nearest,limit=1] item.components."minecraft:custom_data".maf

# dolphins_grace(amp:80) を持つエンティティに対してエフェクト発動
execute as @e[nbt={active_effects:[{id:"minecraft:dolphins_grace",amplifier:80b}]}] at @s run function maf:magic/passive/run_bow_effect with storage maf:tmp arrow_passive

# マーカーエフェクトを除去
effect clear @e[nbt={active_effects:[{id:"minecraft:dolphins_grace",amplifier:80b}]}] minecraft:dolphins_grace

# 半径5以内の最近タグ付き矢を kill
kill @e[type=arrow,tag=maf_passive_arrow,sort=nearest,limit=1,distance=..5]

data remove storage maf:tmp arrow_passive
