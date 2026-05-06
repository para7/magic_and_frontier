function #oh_my_dat:please
execute unless data storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.buff run data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.buff set value []
data remove storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.buff_current
$data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.buff_current set from storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.buff[{buff_category:"$(buff_category)"}]
execute if data storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.buff_current.buff_category if data storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.buff_current.buff_id run function maf:buff/run_destructor with storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.buff_current
data remove storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.buff_current
$data remove storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.buff[{buff_category:"$(buff_category)"}]
$data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.buff_entry set value {buff_id:"$(buff_id)",buff_category:"$(buff_category)",tick:$(tick)}
function maf:buff/add
