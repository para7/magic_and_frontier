tellraw @s [{"text":"設定をクリックで選んでください。\n\n"},{"text":"発光モード: "},{"text":"ON","color":"green","underlined":true,"clickEvent":{"action":"run_command","value":"/data merge storage para7sys:mpu {glow: true}"}},{"text":" / "},{"text":"OFF","color":"red","underlined":true,"clickEvent":{"action":"run_command","value":"/data merge storage para7sys:mpu {glow: false}"}},{"text":"\nベッド使用通知: "},{"text":"ON","color":"green","underlined":true,"clickEvent":{"action":"run_command","value":"/data merge storage para7sys:mpu {bed: true}"}},{"text":" / "},{"text":"OFF","color":"red","underlined":true,"clickEvent":{"action":"run_command","value":"/data merge storage para7sys:mpu {bed: false}"}}]

# execute if data storage para7sys:mpu {glow: true} run data merge storage para7sys:mpu {hoge: false}
# execute if data storage para7sys:mpu {glow: false} run data merge storage para7sys:mpu {hoge: true}

# execute store result storage para7sys:mpu glow byte 1 run data get storage para7sys:mpu hoge

# execute if data storage para7sys:mpu {glow: true} run tellraw @s [{"text":"発光モード"},{"text":"ON", "color":"green"}]
# execute if data storage para7sys:mpu {glow: false} run tellraw @s [{"text":"発光モード"},{"text":"OFF", "color":"red"}]

# effect clear @a minecraft:glowing
