# $(id), $(condition) を macro で受け取る

# condition を tmp に書き出して比較可能にする
$data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.tmp.condition set value "$(condition)"

# none: 自動発動なし → スキップ
execute if data storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.tmp{condition:"none"} run return fail

# attack: 近接攻撃していないときはスキップ
execute if data storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.tmp{condition:"attack"} unless score @s mafMeleeHit matches 1.. run return fail

# always / attack（攻撃あり）: エフェクト実行
$function maf:generated/passive/effect/$(id)
