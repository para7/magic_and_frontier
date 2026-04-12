package export

import (
	"fmt"
	"path/filepath"
	"strings"

	ec "maf_command_editor/app/domain/export/convert"
	bowModel "maf_command_editor/app/domain/model/bow"
)

type BowEffectFunction struct {
	ID   string
	Body string
}

type BowHitFunction struct {
	ID   string
	Body string
}

type BowFlyingFunction struct {
	ID   string
	Body string
}

type BowGroundFunction struct {
	ID   string
	Body string
}

func BuildBowArtifacts(master DBMaster) ([]BowEffectFunction, []BowHitFunction, []BowFlyingFunction, []BowGroundFunction, error) {
	if master == nil {
		return []BowEffectFunction{}, []BowHitFunction{}, []BowFlyingFunction{}, []BowGroundFunction{}, nil
	}

	bows := master.ListBows()
	effects := make([]BowEffectFunction, 0, len(bows))
	hits := make([]BowHitFunction, 0, len(bows))
	flyings := make([]BowFlyingFunction, 0, len(bows))
	grounds := make([]BowGroundFunction, 0, len(bows))

	for _, entry := range bows {
		effects = append(effects, BowEffectFunction{
			ID:   "bow_" + entry.ID,
			Body: buildBowPassiveEffectBody(entry),
		})
		if len(entry.ScriptHit) > 0 {
			hits = append(hits, BowHitFunction{
				ID:   entry.ID,
				Body: strings.Join(entry.ScriptHit, "\n"),
			})
		}
		if len(entry.ScriptFlying) > 0 {
			flyings = append(flyings, BowFlyingFunction{
				ID:   entry.ID,
				Body: strings.Join(entry.ScriptFlying, "\n"),
			})
		}
		if len(entry.ScriptGround) > 0 {
			grounds = append(grounds, BowGroundFunction{
				ID:   entry.ID,
				Body: strings.Join(entry.ScriptGround, "\n"),
			})
		}
	}

	return effects, hits, flyings, grounds, nil
}

func WriteBowArtifacts(effectDir, bowDir, flyingDir, groundDir string, bows []bowModel.BowPassive, effects []BowEffectFunction, hits []BowHitFunction, flyings []BowFlyingFunction, grounds []BowGroundFunction) error {
	for _, entry := range effects {
		path := filepath.Join(effectDir, entry.ID+".mcfunction")
		if err := writeFunctionFile(path, entry.Body); err != nil {
			return err
		}
	}
	for _, entry := range hits {
		path := filepath.Join(bowDir, entry.ID+".mcfunction")
		if err := writeFunctionFile(path, entry.Body); err != nil {
			return err
		}
	}
	for _, entry := range flyings {
		path := filepath.Join(flyingDir, entry.ID+"_flying.mcfunction")
		if err := writeFunctionFile(path, entry.Body); err != nil {
			return err
		}
	}
	for _, entry := range grounds {
		path := filepath.Join(groundDir, entry.ID+"_ground.mcfunction")
		if err := writeFunctionFile(path, entry.Body); err != nil {
			return err
		}
	}
	for _, entry := range bows {
		if len(entry.ScriptHit) == 0 {
			if err := removeFileIfExists(filepath.Join(bowDir, entry.ID+".mcfunction")); err != nil {
				return err
			}
		}
		if len(entry.ScriptFlying) == 0 {
			if err := removeFileIfExists(filepath.Join(flyingDir, entry.ID+"_flying.mcfunction")); err != nil {
				return err
			}
		}
		if len(entry.ScriptGround) == 0 {
			if err := removeFileIfExists(filepath.Join(groundDir, entry.ID+"_ground.mcfunction")); err != nil {
				return err
			}
		}
	}
	return nil
}

func buildBowPassiveEffectBody(entry bowModel.BowPassive) string {
	lifeSub := 1200
	if entry.LifeSub != nil {
		lifeSub = *entry.LifeSub
	}
	lifeValue := 1200 - lifeSub

	lines := []string{
		"execute unless score @s mafBowUsed matches 1.. run return 0",
		"execute store result storage maf:tmp bow_player_id int 1 run scoreboard players get @s mafPlayerID",
		fmt.Sprintf(`execute as @e[type=arrow,distance=..2,nbt=!{inGround:1b},sort=nearest,limit=1] run function maf:bow/tag_bow_arrow {bow_id:%s,life:%d}`, ec.JsonString(entry.ID), lifeValue),
	}
	if len(entry.ScriptHit) > 0 {
		lines = append(lines, "execute as @e[type=arrow,distance=..2,tag=maf_bow_arrow,sort=nearest,limit=1] run function maf:bow/prepare_hit_arrow")
	}
	if len(entry.ScriptFlying) > 0 {
		lines = append(lines, "execute as @e[type=arrow,distance=..2,tag=maf_bow_arrow,sort=nearest,limit=1] run tag @s add flying")
	}
	if len(entry.ScriptGround) > 0 {
		lines = append(lines, "execute as @e[type=arrow,distance=..2,tag=maf_bow_arrow,sort=nearest,limit=1] run tag @s add ground")
	}
	lines = append(lines, entry.ScriptFired...)
	return strings.Join(lines, "\n")
}
