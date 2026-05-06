function #oh_my_dat:please
execute unless data storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.buff[0] run return 0
data remove storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.buff_queue
data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.buff_queue set from storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.buff
data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.buff set value []
function maf:buff/process_queue
data remove storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.buff_queue
data remove storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.buff_current
