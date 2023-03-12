tellraw @s [{"text":"設定をクリックで選んでください。\n\n"},{"text":"発光モード: "},{"text":"ON","color":"green","underlined":true,"clickEvent":{"action":"run_command","value":"/data merge storage p7_:mpu {glow: true}"}},{"text":" / "},{"text":"OFF","color":"red","underlined":true,"clickEvent":{"action":"run_command","value":"/data merge storage p7_:mpu {glow: false}"}},{"text":"\nベッド使用通知: "},{"text":"ON","color":"green","underlined":true,"clickEvent":{"action":"run_command","value":"/data merge storage p7_:mpu {bed: true}"}},{"text":" / "},{"text":"OFF","color":"red","underlined":true,"clickEvent":{"action":"run_command","value":"/data merge storage p7_:mpu {bed: false}"}}]

# execute if data storage p7_:mpu {glow: true} run data merge storage p7_:mpu {hoge: false}
# execute if data storage p7_:mpu {glow: false} run data merge storage p7_:mpu {hoge: true}

# execute store result storage p7_:mpu glow byte 1 run data get storage p7_:mpu hoge

# execute if data storage p7_:mpu {glow: true} run tellraw @s [{"text":"発光モード"},{"text":"ON", "color":"green"}]
# execute if data storage p7_:mpu {glow: false} run tellraw @s [{"text":"発光モード"},{"text":"OFF", "color":"red"}]

# effect clear @a minecraft:glowing
