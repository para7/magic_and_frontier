# ログインした人がいればMP処理を実行
execute if entity @a[scores={mafLogin=1..}] run function maf:system/set_player_id/run

# スコアボードの前処理
execute as @a run function maf:system/score/prescore

# 処理開始
function maf:magic/tick
function maf:soul/tick

# スコアリセット
function maf:system/score/afterscore
