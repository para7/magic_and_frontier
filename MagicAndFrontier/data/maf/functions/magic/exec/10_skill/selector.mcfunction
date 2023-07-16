# tellraw @s [{"text":"スキル【移動速度アップI】をスロット1にセットした"}]

# トリガー有効化
scoreboard players enable @s p7_targetSlot

# フラグ再初期化
scoreboard players set @s p7_targetSlot -1
scoreboard players set @s p7_setSkID -1
scoreboard players set @s p7_setSkEnable 1

# tellraw @s [{"text":"設定をクリックで選んでください。\n\n"},{"text":"スロット1","color":"green","underlined":true,"clickEvent":{"action":"run_command","value":"/data merge storage p7_:mpu {glow: true}"}},{"text":" / "},{"text":"OFF","color":"red","underlined":true,"clickEvent":{"action":"run_command","value":"/data merge storage p7_:mpu {glow: false}"}},{"text":"\nベッド使用通知: "},{"text":"ON","color":"green","underlined":true,"clickEvent":{"action":"run_command","value":"/data merge storage p7_:mpu {bed: true}"}},{"text":" / "},{"text":"OFF","color":"red","underlined":true,"clickEvent":{"action":"run_command","value":"/data merge storage p7_:mpu {bed: false}"}}]
# tellraw @s [{"text":"設定をクリックで選んでください。\n\n"},{"text":"発光モード: "},{"text":"ON","color":"green","underlined":true,"clickEvent":{"action":"run_command","value":"/data merge storage p7_:mpu {glow: true}"}},{"text":" / "},{"text":"OFF","color":"red","underlined":true,"clickEvent":{"action":"run_command","value":"/data merge storage p7_:mpu {glow: false}"}},{"text":"\nベッド使用通知: "},{"text":"ON","color":"green","underlined":true,"clickEvent":{"action":"run_command","value":"/data merge storage p7_:mpu {bed: true}"}},{"text":" / "},{"text":"OFF","color":"red","underlined":true,"clickEvent":{"action":"run_command","value":"/data merge storage p7_:mpu {bed: false}"}}]

# tellraw @s ["",{"text":"隠しジョブ「エンド」","color":"red"},{"text":"の就職条件を満たした！\nマウスを重ねると説明が確認できます: \n"},{"text":"就職条件","color":"aqua","hoverEvent":{"action":"show_text","contents":"エンダーマンを100体以上倒してからエンダードラゴンにトドメを刺す\n(キャンセルした場合、エンダーマンの討伐数は保持されます。)"}},{"text":" ","color":"aqua"},{"text":"ジョブ特性","color":"aqua","hoverEvent":{"action":"show_text","contents":"エンダーパール投擲時にエンダーパールを取得する"}},{"text":" ","color":"aqua"},{"text":"スキル1","color":"aqua","hoverEvent":{"action":"show_text","contents":"ムーブ 向いている方向へ短距離テレポートする"}},{"text":" ","color":"aqua"},{"text":"スキル2","color":"aqua","hoverEvent":{"action":"show_text","contents":"ジャンプ 最寄りの友好モブの場所へワープ"}},{"text":" ","color":"aqua"},{"text":"スキル3","color":"aqua","hoverEvent":{"action":"show_text","contents":"スワップ 最寄りの敵対モブと場所を入れ替える"}},{"text":"\n再度条件を満たすともう一度就職することができます。\n就職しますか?(クリックで決定) : "},{"text":"はい","color":"aqua","clickEvent":{"action":"run_command","value":"/trigger PJEndJobFlag set 1"}},{"text":"     "},{"text":"いいえ","color":"aqua","clickEvent":{"action":"run_command","value":"/trigger PJEndJobFlag set 0"}}]


# TODO: 最大スロット数をスコアボードから取得する
tellraw @s ["",{"text":"設定するスロットを選択してください\n"},{"text":"スロット1","color":"aqua","clickEvent":{"action":"run_command","value":"/trigger p7_targetSlot set 1"}}]