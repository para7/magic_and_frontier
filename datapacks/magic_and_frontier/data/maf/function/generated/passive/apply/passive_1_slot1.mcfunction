data remove storage p7:maf passive.tmp
data modify storage p7:maf passive.tmp.uuid set from entity @s UUID
execute store result storage p7:maf passive.tmp.u0 int 1 run data get storage p7:maf passive.tmp.uuid[0]
execute store result storage p7:maf passive.tmp.u1 int 1 run data get storage p7:maf passive.tmp.uuid[1]
execute store result storage p7:maf passive.tmp.u2 int 1 run data get storage p7:maf passive.tmp.uuid[2]
execute store result storage p7:maf passive.tmp.u3 int 1 run data get storage p7:maf passive.tmp.uuid[3]
data modify storage p7:maf passive.tmp.slot set value 1
data modify storage p7:maf passive.tmp.id set value "passive_1"
function maf:generated/passive/apply/set_slot_by_uuid with storage p7:maf passive.tmp
tellraw @s [{"text":"[slot1]に[Sword Slash]を設定しました"}]
